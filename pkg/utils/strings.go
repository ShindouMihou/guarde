package utils

import (
	"strings"
	"unicode"
)

func HasPrefixStr(source, prefix string) bool {
	return len(source) >= len(prefix) && strings.EqualFold(source[0:len(prefix)], prefix)
}

func Simplify(source string) string {
	return strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII {
			return -1
		}
		if r == '.' || unicode.IsGraphic(r) {
			return r
		}
		if !unicode.IsDigit(r) && !unicode.IsSymbol(r) && !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			return -1
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, source)
}
