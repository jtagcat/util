package util

import "os"

type rollingFile struct {
	currentName func() string
	openFn      func(name string) (*os.File, error)

	file *os.File
}

type OpenFn func(name string) (*os.File, error)

// implements OpenFn
func RollingOpenFn(name string) (*os.File, error) {
	return os.Open(name)
}

// implements OpenFn
func RollingAppendFn(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o660)
}

func NewRollingFile(currentName func() string, openFn OpenFn) rollingFile {
	return rollingFile{
		currentName: currentName,
		openFn:      openFn,
	}
}

func (r *rollingFile) Current() (*os.File, error) {
	newName := r.currentName()
	if newName == r.file.Name() {
		return r.file, nil
	}

	if err := r.Close(); err != nil {
		return nil, err
	}

	f, err := r.openFn(newName)
	r.file = f
	return f, err
}

func (r *rollingFile) Close() error {
	return r.file.Close()
}
