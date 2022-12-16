package flatten_test

import (
	"os"
	"path"
	"testing"

	"github.com/jtagcat/util/flatten"
	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	dir := t.TempDir()
	assert.Nil(t, os.MkdirAll(path.Join(dir, "hello", "world"), 0o770))

	_, err := os.Create(path.Join(dir, "hello", "foo"))
	assert.Nil(t, err)
	_, err = os.Create(path.Join(dir, "hello", "world", "bar"))
	assert.Nil(t, err)

	assert.Nil(t, flatten.Flatten(dir))

	dirS, err := os.ReadDir(dir)
	assert.Nil(t, err)

	type TestDir struct {
		Name  string
		IsDir bool
	}

	var got []TestDir
	for _, f := range dirS {
		got = append(got, TestDir{
			Name:  f.Name(),
			IsDir: f.IsDir(),
		})
	}

	assert.Equal(t, []TestDir{
		{Name: "foo"}, {Name: "bar"},
	}, got)
}

func TestFlattenConflict(t *testing.T) {
	dir := t.TempDir()
	assert.Nil(t, os.MkdirAll(path.Join(dir, "hello", "world"), 0o770))

	_, err := os.Create(path.Join(dir, "hello", "foo"))
	assert.Nil(t, err)
	_, err = os.Create(path.Join(dir, "hello", "world", "foo"))
	assert.Nil(t, err)

	assert.ErrorIs(t, flatten.Flatten(dir), os.ErrExist)
}
