package unicode

import (
	"regexp"
	"unicode"
)

// Tokenizer is a simple unicode tokenizer, that groups characters by code block
// A sequence of characters is treated as one token, as long as they belong to the same unicode code block
// Numerals, spacing and punctuation are treated as separates code blocks
type Tokenizer struct {
	UP             Processor
	SkipWhiteSpace bool
}

type Token struct {
	UnicodeBlock string `json:"block"`
	String       string `json:"string"`
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

const (
	whitespace = "White space"
	numeric    = "Numeric"
	punct      = "Punctuation"
)

// BlockFor returns the name of the unicode block for the input rune. Numerals, spacing and punctuation are treated as separate code blocks.
func (t *Tokenizer) BlockFor(r rune) string {
	n := int(r)
	if numericCharInterval.matches(n) {
		return numeric
	}
	if unicode.IsSpace(r) {
		return whitespace
	}
	if punctRE.MatchString(string(r)) {
		return punct
	}
	return BlockFor(r)
}

func (t *Tokenizer) skippable(token Token) bool {
	if t.SkipWhiteSpace && token.UnicodeBlock == whitespace {
		return true
	}
	return false
}

func (t *Tokenizer) Tokenize(s string) []Token {
	sx := s
	var lastBlock string
	var accu = []rune{}
	var res = []Token{}
	for _, r := range t.UP.Normalize(sx) {
		block := t.BlockFor(r)
		if len(accu) > 0 && block != lastBlock {
			token := Token{String: string(accu), UnicodeBlock: lastBlock}
			if !t.skippable(token) {
				res = append(res, token)
			}
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
