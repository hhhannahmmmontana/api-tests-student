package utils

import (
	"strings"
	"unicode"
)

func MixCase(name string) string {
	var b strings.Builder
	b.Grow(len(name))
	for i, r := range []rune(name) {
		if i%2 == 1 {
			b.WriteRune(unicode.ToUpper(r))
		} else {
			b.WriteRune(unicode.ToLower(r))
		}
	}

	return b.String()
}
