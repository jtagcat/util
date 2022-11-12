package rolling_file

import (
	"io/fs"
	"os"
)

type OpenFn func(name string) (*os.File, error)

// implements OpenFn
func OsOpenFn(name string) (*os.File, error) {
	return os.Open(name)
}

// implements OpenFn
func AppendFn(perms fs.FileMode) OpenFn {
	return func(name string) (*os.File, error) {
		return os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, perms)
	}
}
