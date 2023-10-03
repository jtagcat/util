package std_test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/jtagcat/util/std"
	"github.com/stretchr/testify/assert"
)

func TestRunWithStdouts0(t *testing.T) {
	cmd := exec.Command("bash", "-c", "echo hello stdout; echo hello stderr 1>&2")

	stdout, stderr, err := std.RunWithStdouts(cmd, false)
	assert.Nil(t, err)
	assert.Equal(t, "hello stdout", stdout)
	assert.Equal(t, "hello stderr", stderr)
}

func TestRunWithStdouts1NoEcho(t *testing.T) {
	cmd := exec.Command("bash", "-c", "echo hello stdout; echo hello stderr 1>&2; exit 1")

	stdout, stderr, err := std.RunWithStdouts(cmd, false)
	assert.NotNil(t, err)
	assert.False(t, strings.Contains(err.Error(), "hello"))

	assert.Equal(t, "hello stdout", stdout)
	assert.Equal(t, "hello stderr", stderr)
}

func TestRunWithStdouts1(t *testing.T) {
	cmd := exec.Command("bash", "-c", "echo hello stdout; echo hello stderr 1>&2; exit 1")

	stdout, stderr, err := std.RunWithStdouts(cmd, true)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "hello"))

	assert.Equal(t, "hello stdout", stdout)
	assert.Equal(t, "hello stderr", stderr)
}
