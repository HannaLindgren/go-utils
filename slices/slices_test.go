package slices

import (
	//"fmt"
	"reflect"
	"testing"
)

var fsExpGot = "expected: %#v ; got: %#v"

func TestValidInsert(t *testing.T) {
	var input, expect, result []string

	//  min legal index
	input = []string{"a", "b", "d"}
	expect = []string{"c", "a", "b", "d"}
	result = Insert(input, 0, "c")
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(fsExpGot, expect, result)
	}

	//  legal index
	input = []string{"a", "b", "d"}
	expect = []string{"a", "b", "c", "d"}
	result = Insert(input, 2, "c")
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(fsExpGot, expect, result)
	}

	//  max legal index
	input = []string{"a", "b", "d"}
	expect = []string{"a", "b", "d", "c"}
	result = Insert(input, 3, "c")
	if !reflect.DeepEqual(result, expect) {
		t.Errorf(fsExpGot, expect, result)
	}

}

func TestInvalidInsert1(t *testing.T) {
	var input []string

	defer func() {
		if err := recover(); err != nil {
			//t.Errorf("panic occurred: %v", err)
		} else {
			t.Errorf("expected panic here")
		}
	}()

	//  index too high
	input = []string{"a", "b", "d"}
	Insert(input, 4, "c")

}

func TestInvalidInsert2(t *testing.T) {
	var input []string

	defer func() {
		if err := recover(); err != nil {
			//t.Errorf("panic occurred: %v", err)
		} else {
			t.Errorf("expected panic here")
		}
	}()

	//  negative index
	input = []string{"a", "b", "d"}
	Insert(input, -1, "c")
}
