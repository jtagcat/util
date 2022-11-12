package tail

import (
	"context"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// Unstable, beta
func New(ctx context.Context, name string, offset int64, whence int) (<-chan *Line, <-chan error, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}

	if err := w.Add(name); err != nil {
		return nil, nil, err
	}

	lineChan, errChan := make(chan *Line), make(chan error)

	go singleFile(ctx, w, &Tailable{
		Name: name, Offset: offset, Whence: whence,
	}, lineChan, errChan)

	return lineChan, errChan, nil
}

// assumes file exists
func singleFile(ctx context.Context,
	w *fsnotify.Watcher, file *Tailable,
	lineChan chan<- *Line, errChan chan<- error,
) {
	var wg sync.WaitGroup

	sctx, cancel := context.WithCancel(ctx)
	file.wakeup = make(chan struct{}, 1)

	wg.Add(1)
	go func() {
		// no need to lock/unlock orderedLineChan, as we only have one same-named file across its life
		fileHandle(sctx, *file, true, &orderedLines{c: lineChan}, errChan)
		wg.Done()
	}()

	defer func() {
		w.Close()

		cancel()
		wg.Wait()

		close(lineChan)
		close(errChan)
	}()

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

			switch ev.Op {
			case fsnotify.Write:
				select {
				case file.wakeup <- struct{}{}:
				default:
				}
			case fsnotify.Remove, fsnotify.Rename:
				close(file.wakeup)
				return
			}
		}
	}
}
