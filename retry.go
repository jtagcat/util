package util

import (
	"k8s.io/apimachinery/pkg/util/wait"
)

// Do copy this code instead of depending on this library
//
// similar to "k8s.io/apimachinery/pkg/util/retry"
func RetryOnError(backoff wait.Backoff, fn func() (retryable bool, err error)) error {
	return wait.ExponentialBackoff(backoff, func() (done bool, _ error) {
		retryable, err := fn()
		if err == nil || !retryable {
			return true, err
		}
		return false, nil
	})
}
