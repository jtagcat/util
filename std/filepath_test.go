package std_test

import (
	"fmt"
	"testing"

	"github.com/jtagcat/util/std"
)

func TestSafeJoin(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		base string
		elem []string
		want string
	}{
		{base: "hello", elem: []string{"hell/../world"}, want: "hello/world"},
		{base: "hello", elem: []string{"world/.."}, want: "hello"},
		{base: "/hello/", elem: []string{"../world"}, want: "/hello/world"},
		{base: "hello/", elem: []string{"../world"}, want: "hello/world"},
		{base: "hello", elem: []string{"../world"}, want: "hello/world"},
		{base: "hello", elem: []string{"/../world"}, want: "hello/world"},
		{base: "hello", elem: []string{"../../../../../../../../world"}, want: "hello/world"},
		{base: "hello", elem: []string{"/world"}, want: "hello/world"},
		{base: "hello", elem: []string{"world"}, want: "hello/world"},
	} {
		_ = t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := std.SafeJoin(test.base, test.elem...)

			if got != test.want {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}
