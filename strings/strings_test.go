package strings

import (
	"testing"
)

var fsExpGot = "expected: %#v ; got: %#v"

func TestReverse(t *testing.T) {
	var test = func(in, exp string) {
		got := Reverse(in)
		if got != exp {
			t.Errorf(fsExpGot, got, exp)
		}
	}
	test("anna", "anna")
	test("irevalsittaniallainattislaveri", "irevalsittaniallainattislaveri")
	test("street", "teerts")
	test("Street", "teertS")
}
