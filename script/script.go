package script

import (
	"encoding/json"
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
