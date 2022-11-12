package script

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
)

// golang can be used for scripting :)

// writes object to file as json
func JsonToFile(filename, indent string, object interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	je := json.NewEncoder(f)
	je.SetIndent("", indent)

	return je.Encode(object)
}

// os.Mkdir, if exists return nil (but not os.mkdirAll)
func MkdirExisting(pathstr string) error {
	err := os.Mkdir(pathstr, os.ModePerm)
	if errors.Is(err, fs.ErrExist) {
		return nil
	}
	return err
}
