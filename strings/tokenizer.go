package strings

import (
	"fmt"
	"strings"
	//"os"
	"regexp"
)

// RegexpTokenizer is a simple tokenizer using a regexp as a delimiter definition
type RegexpTokenizer struct {
	delimRE *regexp.Regexp
}

// Split the input string using the specified `delimRE` delimiter definition. Returns a slice of string tokens (including delimiter tokens).
func (t RegexpTokenizer) Split(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	ms := t.delimRE.FindAllStringIndex(s, -1)
	if ms == nil || len(ms) == 0 {
		return []string{s}
	}
	splitPoints := map[int]bool{}
	for _, m := range ms {
		splitPoints[m[0]] = true
		splitPoints[m[1]] = true
	}
	//fmt.Fprintf(os.Stderr, "%v\n", splitPoints)
	res := []string{}
	acc := []rune{}
	for i, r := range s {
		//fmt.Fprintf(os.Stderr, "<%s> %s %v %v\n", string(acc), string(r), i, splitPoints[i])
		if splitPoints[i] {
			if len(acc) > 0 {
				res = append(res, string(acc))
				acc = []rune{}
			}
		}
		acc = append(acc, r)
	}
	if len(acc) > 0 {
		res = append(res, string(acc))
		acc = []rune{}
	}
	test := strings.Join(res, "")
	if test != s {
		msg := fmt.Sprintf("Tokenized string must match input; expected <%s>, got <%s>", s, test)
		panic(msg)
	}
	return res
}
