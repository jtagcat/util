package std_test

import (
	"reflect"
	"testing"
	"unicode"

	"github.com/jtagcat/util/std"
	"github.com/stretchr/testify/assert"
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

func TestStableSplitN(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		s    string
		n    int
		want []string
	}{
		{s: "", n: 2, want: []string{"", ""}},
		{s: "one", n: 2, want: []string{"one", ""}},
		{s: "one,two", n: 2, want: []string{"one", "two"}},
		{s: "one,two,three", n: 2, want: []string{"one", "two,three"}},
		{s: "one,two,three,four", n: 3, want: []string{"one", "two", "three,four"}},
	} {
		_ = t.Run(test.s, func(t *testing.T) {
			got := std.StableSplitN(test.s, ",", test.n)

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("expected %v, got %v", test.want, got)
			}
		})
	}
}

func TestCutFunc(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		s         string
		wantLeft  string
		wantRight string
		wantOk    bool
	}{
		{s: "", wantLeft: "", wantRight: "", wantOk: false},
		{s: "hello", wantLeft: "hello", wantRight: "", wantOk: false},
		{s: "hello123", wantLeft: "hello", wantRight: "23", wantOk: true},
	} {
		_ = t.Run(test.s, func(t *testing.T) {
			leftOf, rightAnd, ok := std.CutFunc(test.s, unicode.IsDigit)

			assert.Equal(t, test.wantLeft, leftOf)
			assert.Equal(t, test.wantRight, rightAnd)
			assert.Equal(t, test.wantOk, ok)
		})
	}
}


func TestRevCutFunc(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		s         string
		wantLeft  string
		wantRight string
		wantOk    bool
	}{
		{s: "", wantLeft: "", wantRight: "", wantOk: false},
		{s: "hello", wantLeft: "hello", wantRight: "", wantOk: false},
		{s: "hellõ", wantLeft: "hellõ", wantRight: "", wantOk: false},
		{s: "hello123", wantLeft: "hello", wantRight: "23", wantOk: true},
		{s: "hellõ123", wantLeft: "hellõ", wantRight: "23", wantOk: true},
	} {
		_ = t.Run(test.s, func(t *testing.T) {
			leftOf, rightAnd, ok := std.CutFunc(test.s, unicode.IsDigit)

			assert.Equal(t, test.wantLeft, leftOf)
			assert.Equal(t, test.wantRight, rightAnd)
			assert.Equal(t, test.wantOk, ok)
		})
	}
}
