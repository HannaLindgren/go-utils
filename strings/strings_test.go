package strings

import (
	"testing"
)

var fsExpGot = "expected: %#v ; got: %#v"

func TestReverse(t *testing.T) {
	var test = func(in, exp string) {
		got := Reverse(in)
		if got != exp {
			t.Errorf(fsExpGot, exp, got)
		}
	}
	test("anna", "anna")
	test("irevalsittaniallainattislaveri", "irevalsittaniallainattislaveri")
	test("street", "teerts")
	test("Street", "teertS")
}

func TestUpcaseInitial(t *testing.T) {
	var test = func(in, exp string, downcaseRest bool) {
		got := UpcaseInitial(in, downcaseRest)
		if got != exp {
			t.Errorf(fsExpGot, exp, got)
		}
	}
	test("annA", "Anna", true)
	test("annA", "AnnA", false)
	test("street Feet", "Street Feet", false)
	test("street Feet", "Street feet", true)
}
