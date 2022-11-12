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

func TestTailFileAddNothing(t *testing.T) {
	t.Parallel()
	ctx, name := context.Background(), filepath.Join(t.TempDir(), "foo")
	assert.Nil(t,
		os.WriteFile(name, nil, 0o600))

	lines, errs, err := tail.New(ctx, name, 0, io.SeekStart)
	assert.Nil(t, err)

	assert.Nil(t,
		os.WriteFile(name, nil, 0o600))

	select {
	case e := <-errs:
		t.Errorf("did not expect any error, got: %e", e)
	case l := <-lines:
		t.Errorf("did not expect any lines, got: %v", l)
	default:
	}
}
