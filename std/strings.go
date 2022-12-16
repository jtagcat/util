package std

import (
	"strings"
)

// strings.Cut, but starting from last character, found is either empty or seperator
func RevCut(s, sep string) (leftOf, rightOf string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func TrimLen(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}

	return string(r[:max])
}
