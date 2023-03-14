package strings

import (
	"strings"
)

// UpcaseInitial Upcase the first rune of a string, and optionally downcase the remainder
func UpcaseInitial(s string, downcaseRemainder bool) string {
	runes := []rune(s)
	head := ""
	if len(runes) > 0 {
		head = strings.ToUpper(string(runes[0]))
	}
	tail := ""
	if len(runes) > 0 {
		tail = string(runes[1:])
		if downcaseRemainder {
			tail = strings.ToLower(tail)
		}
	}
	return head + tail
}

// Reverse a string
func Reverse(str string) string {
	bts := []rune(str)
	for i, j := 0, len(bts)-1; i < j; i, j = i+1, j-1 {
		bts[i], bts[j] = bts[j], bts[i]
	}
	return string(bts)
}
