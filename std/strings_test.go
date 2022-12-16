package std_test

import (
	"testing"

	"github.com/jtagcat/util/std"
)

func TestTrimLen(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		s    string
		max  int
		want string
	}{
		{s: "aõ", max: 2, want: "aõ"},
		{s: "aõö", max: 2, want: "aõ"},
	} {
		_ = t.Run(test.s, func(t *testing.T) {
			got := std.TrimLen(test.s, test.max)
			if got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
}
