package simple

import (
	"time"
)

// Deprecated: copy this code instead of depending on this library
//
// max and maxWait must be used
func Batch[T any](max int, maxWait time.Duration,
	stream <-chan T, batched chan<- []T,
) {
	if max < 2 {
		panic("simple.Batch: max must be at least 2")
	}
	if maxWait <= 0 {
		panic("simple.Batch: maxWait must be positive")
	}

	for {
		var batch []T
		send := func() {
			if len(batch) > 0 {
				batched <- batch
			}
		}

		timer := time.NewTimer(maxWait)

	collector:
		for i := 0; i < max; i++ {
			select {
			case s, ok := <-stream:
				if ok {
					batch = append(batch, s)
					continue
				}

				send()
				close(batched)
				return

			case <-timer.C:
				break collector
			}
		}

		send()
	}
}
