package util

import (
	"encoding/csv"
	"io/fs"
	"os"
)

type rollingFile struct {
	currentName func() string
	openFn      func(name string) (*os.File, error)

	preClose func() error
	file     *os.File
}

type OpenFn func(name string) (*os.File, error)

// implements OpenFn
func RollingOpenFn(name string) (*os.File, error) {
	return os.Open(name)
}

// implements OpenFn
func RollingAppendFn(perms fs.FileMode) OpenFn {
	return func(name string) (*os.File, error) {
		return os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, perms)
	}
}

func NewRollingFile(currentName func() string, openFn OpenFn) rollingFile {
	return rollingFile{
		currentName: currentName,
		openFn:      openFn,
	}
}

func (r *rollingFile) Current() (_ *os.File, _ error, changed bool) {
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

func (r *rollingFile) Close() error {
	if err := r.preClose(); err != nil {
		return err
	}

	return r.file.Close()
}

type rollingCsvAppender struct {
	rollingFile

	csv *csv.Writer
}

func NewRollingCsvAppender(currentName func() string, perms fs.FileMode) rollingCsvAppender {
	a := rollingCsvAppender{
		rollingFile: rollingFile{
			currentName: currentName,
			openFn:      RollingAppendFn(perms),
		},
	}

	a.preClose = func() error {
		if a.csv != nil {
			a.csv.Flush()
		}
		return nil
	}

	return a
}

func (c *rollingCsvAppender) Current() (_ *csv.Writer, _ error, changed bool) {
	f, err, changed := c.rollingFile.Current()

	if err == nil && changed {
		c.csv = csv.NewWriter(f)
	}

	return c.csv, err, changed
}

// Writes are buffered, so Flush must eventually be called to ensure
// that the record is written to the underlying io.Writer.
func (c *rollingCsvAppender) WriteCurrent(record []string) error {
	writer, err, _ := c.Current()
	if err != nil {
		return err
	}

	return writer.Write(record)
}
