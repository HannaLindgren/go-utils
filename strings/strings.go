package strings

import (
	"fmt"
	"strings"
)

// UpcaseInitial Upcase the first rune of a string, and downcase the remainder
func UpcaseInitial(s string) string {
	runes := []rune(s)
	head := ""
	if len(runes) > 0 {
		head = strings.ToUpper(string(runes[0]))
	}
	tail := ""
	if len(runes) > 0 {
		tail = strings.ToLower(string(runes[1:]))
	}
	return head + tail
}

var delims = strings.Split(" 	(){}\"”<>'-/&;.!?", "")

func simpleTokenize(s string) []string {
	var acc = []string{s}
	for _, delim := range delims {
		tmp := []string{}
		for _, s := range acc {
			tmp = append(tmp, strings.SplitAfter(s, delim)...)
		}
		acc = tmp
	}
	return acc
}

// CapitalizeTokens Capitalize each token of a string (using a rudimentary tokenization)
func CapitalizeTokens(s string) string {
	res := ""
	for _, token := range simpleTokenize(s) {
		res = res + UpcaseInitial(token)
	}
	if !strings.EqualFold(s, res) {
		panic(fmt.Sprintf("Expected output string to equal input string except for case, but found: <%s> => <%s>", s, res))
	}
	return res
}
