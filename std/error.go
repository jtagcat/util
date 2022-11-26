package std

import "errors"

// for errors.Is(err, ERr)
type GenericErr struct {
	Err     error
	Wrapped error
}

func (a GenericErr) Is(target error) bool {
	return errors.Is(a.Err, target)
}

func (a GenericErr) Unwrap() error {
	return a.Wrapped
}

func (a GenericErr) Error() string {
	return a.Err.Error() + ": " + a.Wrapped.Error()
}
