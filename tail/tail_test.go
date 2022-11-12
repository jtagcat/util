package tail_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/jtagcat/util/tail"
	"github.com/stretchr/testify/assert"
)

// TODO: symlink change
// TODO: write with file not changing
// TODO: replace file

// expected: fail
func TestExampleSingle(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file, file2 := filepath.Join(dir, "foo"), filepath.Join(dir, "bar")

	assert.Nil(t,
		os.WriteFile(file, []byte("hello"), 0o600))

	l, e, err := tail.New(context.Background(), file, 0, io.SeekStart)
	assert.Nil(t, err)

	// first empty
	select {
	case <-e:
		t.FailNow()
	case s := <-l:
		assert.Equal(t, file, *s.Filename)
		assert.Equal(t, "hello", string(s.Bytes))
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
	assert.Equal(t, "world", string(s.Bytes))
	s = <-l
	assert.Equal(t, "space", string(s.Bytes))
	select {
	case <-e:
		t.FailNow()
	case s := <-l:
		assert.Equal(t, len(s.Bytes), 0)
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
		assert.Equal(t, "bar", string(s.Bytes))
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
		assert.Equal(t, "two", string(s.Bytes))
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
	defer f.Close()

	if _, err := f.Write([]byte(s)); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
