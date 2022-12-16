package retry

import (
	"context"

	"k8s.io/apimachinery/pkg/util/wait"
)

// Do copy this code instead of depending on this library
//
// similar to "k8s.io/apimachinery/pkg/util/retry"
func OnError(backoff wait.Backoff, fn func() (retryable bool, err error)) error {
	return wait.ExponentialBackoff(backoff, func() (done bool, _ error) {
		retryable, err := fn()
		if err == nil || !retryable {
			return true, nil
		}
		return false, nil
	})
}

// error is only returned by context
func OnErrorManagedBackoff(ctx context.Context, backoff wait.Backoff, fn func() (retryable bool, err error)) error {
	return wait.ManagedExponentialBackoffWithContext(ctx, backoff, func() (done bool, _ error) {
		retryable, err := fn()
		if err == nil || !retryable {
			return true, nil
		}
		return false, nil
	})
}
