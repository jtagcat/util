package util

import "os"

type rollingFile struct {
	currentName func() string
	openFn      func(name string) (*os.File, error)

	file *os.File
}

func NewRollingFile(currentName func() string, openFn *func(name string) (*os.File, error)) rollingFile {
	fn := func(name string) (*os.File, error) {
		return os.Open(name)
	}

	if openFn != nil {
		fn = *openFn
	}

	return rollingFile{
		currentName: currentName,
		openFn:      fn,
	}
}

func (r *rollingFile) Current() (*os.File, error) {
	newName := r.currentName()
	if newName == r.file.Name() {
		return r.file, nil
	}

	f, err := r.openFn(newName)
	r.file = f
	return f, err
}
