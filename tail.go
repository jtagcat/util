package util

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Line struct {
	Filename   *string
	Bytes      []byte
	EndOffset  int64 // io.SeekStart
	ReachedEOF bool
}

// File starts tailing file from offset and whence (os.Seek()).
// It follows target file for appends, truncations, and replacements.
// Errors abort connected operations.

type Tailable struct {
	Name string
	// os.Seek() on first open:
	Offset int64
	Whence int

	// loop-persisting
	existed bool

	// copied use
	wakeup chan struct{}
}

var ErrScatteredFiles = errors.New("all Tailable files must be in the same directory")

// Unstable, beta
//
// All files must be in the same directory.
// Channels will be closed after file is deleted //TODO:
func TailFiles(ctx context.Context, files []Tailable) (<-chan *Line, <-chan error, error) {
	if len(files) == 0 {
		return nil, nil, nil
	}

	parentDir := filepath.Dir(files[0].Name)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}

	for i := range files {
		file := &files[i]
		if filepath.Dir(file.Name) != parentDir {
			return nil, nil, ErrScatteredFiles
		}

		// Simulate Create events for files already existing
		_, err = os.Stat(file.Name)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, nil, err
		}

		file.existed = true

		name := file.Name
		go func() {
			w.Events <- fsnotify.Event{
				Name: name,
				Op:   fsnotify.Create,
			}
		}()
	}

	if err := w.Add(parentDir); err != nil {
		return nil, nil, fmt.Errorf("watching parent directory %q: %w", parentDir, err)
	}

	lineChan, errChan := make(chan *Line), make(chan error)

	go tailFiles(ctx, w, &files, lineChan, errChan)

	return lineChan, errChan, err
}

// Consumes Watcher
func tailFiles(ctx context.Context,
	w *fsnotify.Watcher, files *[]Tailable,
	lineChan chan<- *Line, errChan chan<- error,
) {
	defer func() {
		close(lineChan)
		close(errChan)
	}()

	type mapWrap struct {
		*Tailable
		seen     bool
		lineChan orderedLineChan
	}

	names := make(map[string]*mapWrap)
	for _, file := range *files {

		c := make(chan *Line)
		names[filepath.Base(file.Name)] = &mapWrap{
			Tailable: &file,
			lineChan: orderedLineChan{c: c},
		}

		go func() { // relay files
			for {
				l, ok := <-c
				if !ok {
					return
				}
				lineChan <- l
			}
		}()
	}

	for {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}
		// There is no priority in select.
		// To prevent race of looping with expired context,
		// ctx.Done() is checked again at the start.
		select {
		case <-ctx.Done():
			continue
		case err, ok := <-w.Errors:
			if ok {
				errChan <- err
			}
			return

		case ev, ok := <-w.Events:
			if !ok {
				return
			}

			file, ok := names[filepath.Base(ev.Name)]
			if !ok {
				continue
			}

			switch ev.Op {
			case fsnotify.Create:
				var isFakeCreate bool

				if file.seen {
					close(file.wakeup)
				} else {
					file.seen = true

					if file.existed {
						isFakeCreate = true
					}
				}

				file.wakeup = make(chan struct{}, 1)

				// do not give pointer to file, as multiple FDs with same name may exist, with different wakeups
				go func() {
					file.lineChan.Lock()
					fileHandle(ctx, *file.Tailable, isFakeCreate, &file.lineChan, errChan)
					file.lineChan.Unlock()
				}()

			case fsnotify.Write:
				select {
				case file.wakeup <- struct{}{}:
				default:
				}
			}
		}
	}
}

type orderedLineChan struct {
	c          chan<- *Line
	sync.Mutex // guarantee ordered lines across multiple files of same name
}

// Handles tailing a file within its lifespan
func fileHandle(ctx context.Context, file Tailable, useOffset bool, lineChan *orderedLineChan, errChan chan<- error) {
	f, err := os.Open(file.Name)
	if err != nil && !errors.Is(err, os.ErrNotExist) { // ignore ErrNotExist, as it may have been race deleted
		errChan <- err
		return
	}
	defer f.Close()

	var offset int64
	if useOffset {
		offset, err = f.Seek(file.Offset, file.Whence)
		if err != nil {
			errChan <- err
			return
		}
	}

	first, breakNext, b := true, false, bufio.NewReader(f)
	for {

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		if first {
			first = false
		} else {
			offset, err = detectTruncation(f, offset)
			if err != nil {
				errChan <- err
				return
			}
		}

		offset, err = readToEOF(b, &file.Name, offset, lineChan.c)
		if err != nil {
			errChan <- err
			return
		}

		if breakNext {
			return
		}
		select {
		case <-ctx.Done():
		case _, ok := <-file.wakeup:
			if !ok {
				fs, err := f.Stat()
				if err != nil && !errors.Is(err, os.ErrNotExist) {
					errChan <- ctx.Err()
					return
				}

				if err != nil || fs.Size() == offset {
					return
				}

				breakNext = true
			}
		}
	}
}

// offset is io.SeekStart
func detectTruncation(f *os.File, offset int64) (int64, error) {
	fs, err := f.Stat()
	if err != nil && !errors.Is(err, os.ErrNotExist) { // ignore ErrNotExist, as it may have been race deleted
		return 0, err
	}

	if fs.Size() < offset {
		// file has been truncated
		return f.Seek(0, io.SeekStart)
	}

	return offset, nil
}

func readToEOF(buf *bufio.Reader, name *string, offset int64, c chan<- *Line) (int64, error) {
	for {
		b, err := buf.ReadBytes('\n')
		offset += int64(len(b))

		if err != nil && !errors.Is(err, io.EOF) {
			return offset, err
		}

		if err == nil {
			b = b[:len(b)-1] // remove \n
		}

		c <- &Line{
			Filename:   name,
			Bytes:      b,
			EndOffset:  offset,
			ReachedEOF: err != nil,
		}

		if err != nil { // EOF
			return offset, nil
		}
	}
}
