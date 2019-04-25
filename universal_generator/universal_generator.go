package universal_generator

import (
	"strings"
)

// TYPES AND UTILITIES

type input []string
type template []input
type generated []string

func (g generated) String() string {
	var res string
	res = strings.Join(g, " ")
	res = strings.ReplaceAll(res, "  ", " ")
	res = strings.TrimSpace(res)
	return res
}

// EXPAND TEMPLATES

func copyOf(input []string) []string {
	res := []string{}
	for _, s := range input {
		res = append(res, s)
	}
	return res
}

func expandLoop(head input, tail []input, accs []generated) []generated {
	res := []generated{}
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

func expand(template []input) []generated {
	return expandLoop(template[0], template[1:], []generated{{}})
}
