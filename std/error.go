package std

import (
	"errors"

	"golang.org/x/exp/slog"
)

func SlogErr(err error) slog.Attr {
	str := ""
	if err != nil {
		str = err.Error()
	}

	return slog.String("err", str)
}

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
