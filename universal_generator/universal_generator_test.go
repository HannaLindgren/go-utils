package universal_generator

import (
	"reflect"
	"testing"
)

var fsExpGot = "expected: %#v ; got: %#v"

func Test1(t *testing.T) {
	template := Template{
		Input{"idag", "imorgon", "på fredag", "nästa vecka"},
		Input{"ska", "kan", "kommer"},
		Input{"det"},
		Input{"inte", "kanske", ""},
		Input{"regna", "snöa", "brinna", "hagla", "dugga", "blåsa", "storma", "blåsa upp till orkan", "vara uppehåll", "vara fint väder", "vara hög brandrisk i alla län"},
	}
	result := len(template.Expand())
	expect := 1
	for _, i := range template {
		expect = expect * len(i)
	}

	if result != expect {
		t.Errorf(fsExpGot, expect, result)
	}
	if result != 396 {
		t.Errorf(fsExpGot, expect, 396)
	}
}

func Test2(t *testing.T) {
	template := Template{
		Input{"idag", "imorgon", "nästa vecka"},
		Input{"ska", "kan", "kommer"},
		Input{"det"},
		Input{"inte", ""},
		Input{"regna", "snöa"},
	}
	result := template.Expand()
	expect := []Output{
		[]string{"idag", "ska", "det", "inte", "regna"},
		[]string{"idag", "ska", "det", "inte", "snöa"},
		[]string{"idag", "ska", "det", "", "regna"},
		[]string{"idag", "ska", "det", "", "snöa"},
		[]string{"idag", "kan", "det", "inte", "regna"},
		[]string{"idag", "kan", "det", "inte", "snöa"},
		[]string{"idag", "kan", "det", "", "regna"},
		[]string{"idag", "kan", "det", "", "snöa"},
		[]string{"idag", "kommer", "det", "inte", "regna"},
		[]string{"idag", "kommer", "det", "inte", "snöa"},
		[]string{"idag", "kommer", "det", "", "regna"},
		[]string{"idag", "kommer", "det", "", "snöa"},
		[]string{"imorgon", "ska", "det", "inte", "regna"},
		[]string{"imorgon", "ska", "det", "inte", "snöa"},
		[]string{"imorgon", "ska", "det", "", "regna"},
		[]string{"imorgon", "ska", "det", "", "snöa"},
		[]string{"imorgon", "kan", "det", "inte", "regna"},
		[]string{"imorgon", "kan", "det", "inte", "snöa"},
		[]string{"imorgon", "kan", "det", "", "regna"},
		[]string{"imorgon", "kan", "det", "", "snöa"},
		[]string{"imorgon", "kommer", "det", "inte", "regna"},
		[]string{"imorgon", "kommer", "det", "inte", "snöa"},
		[]string{"imorgon", "kommer", "det", "", "regna"},
		[]string{"imorgon", "kommer", "det", "", "snöa"},
		[]string{"nästa vecka", "ska", "det", "inte", "regna"},
		[]string{"nästa vecka", "ska", "det", "inte", "snöa"},
		[]string{"nästa vecka", "ska", "det", "", "regna"},
		[]string{"nästa vecka", "ska", "det", "", "snöa"},
		[]string{"nästa vecka", "kan", "det", "inte", "regna"},
		[]string{"nästa vecka", "kan", "det", "inte", "snöa"},
		[]string{"nästa vecka", "kan", "det", "", "regna"},
		[]string{"nästa vecka", "kan", "det", "", "snöa"},
		[]string{"nästa vecka", "kommer", "det", "inte", "regna"},
		[]string{"nästa vecka", "kommer", "det", "inte", "snöa"},
		[]string{"nästa vecka", "kommer", "det", "", "regna"},
		[]string{"nästa vecka", "kommer", "det", "", "snöa"},
	}

	if !reflect.DeepEqual(expect, result) {
		t.Errorf(fsExpGot, expect, result)
	}
}
