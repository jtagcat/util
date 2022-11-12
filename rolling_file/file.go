package rolling_file

import (
	"os"
)

type file struct {
	currentName func() string
	openFn      func(name string) (*os.File, error)

	preClose func() error
	file     *os.File
}

func New(currentName func() string, openFn OpenFn) file {
	return file{
		currentName: currentName,
		openFn:      openFn,
	}
}

func (r *file) Current() (_ *os.File, _ error, changed bool) {
	newName := r.currentName()

	if r.file != nil {
		if newName == r.file.Name() {
			return r.file, nil, false
		}

		if err := r.preClose(); err != nil {
			return nil, err, true
		}

		if err := r.Close(); err != nil {
			return nil, err, true
		}
	}

	f, err := r.openFn(newName)
	r.file = f
	return f, err, true
}

func (r *file) Close() error {
	if err := r.preClose(); err != nil {
		return err
	}

	return r.file.Close()
}
