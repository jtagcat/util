package tail

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// File starts tailing file from offset and whence (os.Seek()).
// It follows target file for appends, truncations, and replacements.
// Errors abort connected operations.

var ErrScatteredFiles = errors.New("all Tailable files must be in the same directory")

// Unstable, beta
//
// All files must be in the same directory.
// Channels will be closed after file is deleted //TODO:
func Files(ctx context.Context, files []Tailable) (<-chan *Line, <-chan error, error) {
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

	go multipleFiles(ctx, w, &files, lineChan, errChan)

	return lineChan, errChan, err
}

// Consumes Watcher
func multipleFiles(ctx context.Context,
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
		lineChan orderedLines
	}

	names := make(map[string]*mapWrap)
	for _, file := range *files {

		c := make(chan *Line)
		names[filepath.Base(file.Name)] = &mapWrap{
			Tailable: &file,
			lineChan: orderedLines{c: c},
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
