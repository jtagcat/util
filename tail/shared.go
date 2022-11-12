package tail

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"sync"
)

type Line struct {
	Filename   *string
	Bytes      []byte
	EndOffset  int64 // io.SeekStart
	ReachedEOF bool
}

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

type orderedLines struct {
	c          chan<- *Line
	sync.Mutex // guarantee ordered lines across multiple files of same name
}

// Handles tailing a file within its lifespan
func fileHandle(ctx context.Context, file Tailable, useOffset bool, lineChan *orderedLines, errChan chan<- error) {
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
