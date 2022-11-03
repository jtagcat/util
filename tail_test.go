package util_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	util "github.com/jtagcat/util"
	"github.com/stretchr/testify/assert"
)

// TODO: symlink change

func TestExampleSingle(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file, file2 := filepath.Join(dir, "foo"), filepath.Join(dir, "bar")

	assert.Nil(t,
		os.WriteFile(file, []byte("hello"), 0o600))

	l, e, err := util.TailFile(context.Background(), file, 0, io.SeekStart)
	assert.Nil(t, err)

	// first empty
	select {
	case <-e:
		t.FailNow()
	case s := <-l:
		assert.Equal(t, file, *s.Filename)
		assert.Equal(t, "hello", s.String)
	}

	// empty chan
E:
	for {
		select {
		case <-l:
			t.FailNow()
		default:
			break E
		}
	}

	assert.Nil(t,
		append(file, "world\nspace\n"))

	// last double empty

	s := <-l
	assert.Equal(t, "world", s.String)
	s = <-l
	assert.Equal(t, "space", s.String)
	select {
	case <-e:
		t.FailNow()
	case s := <-l:
		assert.Equal(t, "", s.String)
		assert.Equal(t, file, *s.Filename)
	}

	// truncate
	assert.Nil(t,
		os.Truncate(file, 0))
	assert.Nil(t,
		append(file, "bar"))

	select {
	case <-e:
		t.FailNow()
	case s := <-l:
		assert.Equal(t, "bar", s.String)
		assert.Equal(t, file, *s.Filename)
	}

	// replace
	assert.Nil(t,
		os.WriteFile(file2, []byte("two"), 0o600))
	assert.Nil(t,
		os.Rename(file2, file))

	select {
	case <-e:
		t.FailNow()
	case s := <-l:
		assert.Equal(t, "two", s.String)
		assert.Equal(t, file, *s.Filename)
	}

	// empty chan
F:
	for {
		select {
		case <-l:
			t.FailNow()
		default:
			break F
		}
	}
}

func append(name, s string) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}

	if _, err := f.Write([]byte(s)); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
