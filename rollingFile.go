package util

import "os"

type rollingFile struct {
	currentName func() string

	file *os.File
}

func NewRollingFile(currentName func() string) rollingFile {
	return rollingFile{
		currentName: currentName,
	}
}

func (r *rollingFile) Current() (*os.File, error) {
	newName := r.currentName()
	if newName == r.file.Name() {
		return r.file, nil
	}

	f, err := os.Open(newName)
	r.file = f
	return f, err
}
