package wakeup

import (
	"context"
	"errors"
	"fmt"
)

type wakeup chan struct{}

const contextKey = "wakeup"

// WithWakeup returns a copy of parent with an embedded wakeup channel.
//
//	wctx, wake := wakeup.WithWakeup(context.Background())
//	go wakeup.Wait(wctx, func(ctx context.Context) (goToSleep bool) {
//		// do stuff
//	})
//
//	wake.Wakeup()
func WithWakeup(parent context.Context) (context.Context, wakeup) {
	wake := make(wakeup, 1)

	return context.WithValue(context.Background(), contextKey, wake), wake
}

func (w wakeup) Wakeup() {
	select {
	case w <- struct{}{}:
	default:
	}
}

var errNoWakeup = errors.New(fmt.Sprintf("context does not have wakeup (as value %s)", contextKey))

// Wait runs function until goToSleep is true, then waits until Wakeup() or context.Done().
func Wait(wctx context.Context, fn func(ctx context.Context) (goToSleep bool)) error {
	wake, ok := wctx.Value(contextKey).(wakeup)
	if !ok {
		return errNoWakeup
	}

	select {
	case <-wctx.Done():
		return wctx.Err()
	default:
	}

	for {

	run:
		for {
			if sleep := fn(wctx); sleep {
				break run
			}
		}

		select {
		case <-wctx.Done():
			return wctx.Err()

		case _, ok := <-wake:
			if !ok {
				return context.Canceled
			}
		}
	}
}
