package util

import (
	"encoding/csv"
	"os"
)

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

func (r *rollingFile) Current() (_ *os.File, _ error, changed bool) {
	newName := r.currentName()
	if newName == r.file.Name() {
		return r.file, nil, false
	}

	if err := r.Close(); err != nil {
		return nil, err, true
	}

	f, err := r.openFn(newName)
	r.file = f
	return f, err, true
}

func (r *rollingFile) Close() error {
	return r.file.Close()
}

type rollingCsvAppender struct {
	rollingFile

	csv *csv.Writer
}

func NewRollingCsvAppender(currentName func() string) rollingCsvAppender {
	return rollingCsvAppender{
		rollingFile: rollingFile{
			currentName: currentName,
			openFn:      RollingAppendFn,
		},
	}
}

func (c *rollingCsvAppender) Current() (_ *csv.Writer, _ error, changed bool) {
	f, err, changed := c.rollingFile.Current()

	if err == nil && changed {
		c.csv = csv.NewWriter(f)
	}

	return c.csv, err, changed
}
