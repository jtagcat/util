package rolling_file

import (
	"encoding/csv"
	"io/fs"
)

type csvAppender struct {
	file

	csv *csv.Writer
}

func NewCsvAppender(currentName func() string, perms fs.FileMode) csvAppender {
	a := csvAppender{
		file: file{
			currentName: currentName,
			openFn:      AppendFn(perms),
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

func (c *csvAppender) Current() (_ *csv.Writer, _ error, changed bool) {
	f, err, changed := c.file.Current()

	if err == nil && changed {
		c.csv = csv.NewWriter(f)
	}

	return c.csv, err, changed
}

// Writes are buffered, so Flush must eventually be called to ensure
// that the record is written to the underlying io.Writer.
func (c *csvAppender) WriteCurrent(record []string) error {
	writer, err, _ := c.Current()
	if err != nil {
		return err
	}

	return writer.Write(record)
}
