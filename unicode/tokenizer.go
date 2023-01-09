package unicode

import (
	"regexp"
)

// Tokenizer is a simple unicode tokenizer, that groups characters by code block
// A sequence of characters is treated as one token, as long as they belong to the same unicode code block
// Numerals, spacing and punctuation are treated as separates code blocks
type Tokenizer struct {
	UP Processor
}

type Token struct {
	UnicodeBlock string
	String       string
}

type interval struct {
	from int
	to   int
}

func (i interval) matches(n int) bool {
	return n >= i.from && n <= i.to
}

var numericCharInterval = interval{48, 57}
var spacingRE = regexp.MustCompile(`\s`)
var punctRE = regexp.MustCompile(`\pP`)

// BlockFor returns the name of the unicode block for the input rune. Numerals, spacing and punctuation are treated as separate code blocks.
func (t *Tokenizer) BlockFor(r rune) string {
	n := int(r)
	if numericCharInterval.matches(n) {
		return "Numeric"
	}
	if spacingRE.MatchString(string(r)) {
		return "White space"
	}
	if punctRE.MatchString(string(r)) {
		return "Punctuation"
	}
	return BlockFor(r)
}

func (t *Tokenizer) Tokenize(s string) []Token {
	sx := s
	var lastBlock string
	var accu = []rune{}
	var res = []Token{}
	for _, r := range t.UP.Normalize(sx) {
		block := BlockFor(r)
		if len(accu) > 0 && block != lastBlock {
			res = append(res, Token{String: string(accu), UnicodeBlock: block})
			accu = []rune{}
		}
		accu = append(accu, r)
		lastBlock = block
	}
	if len(accu) > 0 {
		res = append(res, Token{String: string(accu), UnicodeBlock: lastBlock})
	}
	return res
}