package universal_generator

import (
	"strings"
)

// TYPES AND UTILITIES

// Input is a set of interchangeable words, phonemes, or other units
type Input []string

// Template is a sequence of Inputs that can form a complete utterance or string or similar.
// Usage examples in universal_generator_test.go
type Template []Input

// Output is a sequence of words, phonemes, units, generated from a template. From one template, you will get a slice of Outputs.
type Output []string

// String generates an output string with all elements joined by white space
func (o Output) String() string {
	return o.Join(" ")
}

// Join the elements with the specified separator string
func (o Output) Join(separator string) string {
	res := []string{}
	for _, s := range o {
		if s != "" {
			res = append(res, s)
		}
	}
	return strings.TrimSpace(strings.Join(res, separator))
}

func copyOf(Input []string) []string {
	res := []string{}
	res = append(res, Input...)
	return res
}

// EXPANSION ALGORITHM

func expandLoop(head Input, tail []Input, accs []Output) []Output {
	res := []Output{}
	for _, acc := range accs {
		for _, add := range head {
			newAcc := append(copyOf(acc), add)
			res = append(res, newAcc)
		}
	}
	if len(tail) == 0 {
		return res
	}
	return expandLoop(tail[0], tail[1:], res)
}

// Expand the input template
func (t Template) Expand() []Output {
	return expandLoop(t[0], t[1:], []Output{{}})
}
