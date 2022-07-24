package jsonschema

import (
	"strconv"
	"strings"
	"unicode"
)

func trimString(s string) string {
	return strings.TrimFunc(s, func(r rune) bool { return unicode.IsSpace(r) })
}

func reverseStringSlice(s []string) []string {
	for left, right := 0, len(s)-1; left < right; left, right = left+1, right-1 {
		s[left], s[right] = s[right], s[left]
	}
	return s
}

func indexString(i int) string {
	return "[" + strconv.FormatInt(int64(i), 10) + "]"
}
