package universal_generator

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

var fsExpGot = "expected: %#v ; got: %#v"

func f0() {
	fmt.Fprintf(os.Stderr, "")
}

func Test1(t *testing.T) {
	template := template{
		input{"idag", "imorgon", "på fredag", "nästa vecka"},
		input{"ska", "kan", "kommer"},
		input{"det"},
		input{"inte", "kanske", ""},
		input{"regna", "snöa", "brinna", "hagla", "dugga", "blåsa", "storma", "blåsa upp till orkan", "vara uppehåll", "vara fint väder", "vara hög brandrisk i alla län"},
	}
	result := len(expand(template))
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
	template := template{
		input{"idag", "imorgon"},
		input{"ska", "kan", "kommer"},
		input{"det"},
		input{"inte", ""},
		input{"regna", "snöa"},
	}
	result := expand(template)
	expect := []output{
		strings.Split("idag ska det inte regna", " "),
		strings.Split("idag ska det inte snöa", " "),
		strings.Split("idag ska det  regna", " "),
		strings.Split("idag ska det  snöa", " "),
		strings.Split("idag kan det inte regna", " "),
		strings.Split("idag kan det inte snöa", " "),
		strings.Split("idag kan det  regna", " "),
		strings.Split("idag kan det  snöa", " "),
		strings.Split("idag kommer det inte regna", " "),
		strings.Split("idag kommer det inte snöa", " "),
		strings.Split("idag kommer det  regna", " "),
		strings.Split("idag kommer det  snöa", " "),
		strings.Split("imorgon ska det inte regna", " "),
		strings.Split("imorgon ska det inte snöa", " "),
		strings.Split("imorgon ska det  regna", " "),
		strings.Split("imorgon ska det  snöa", " "),
		strings.Split("imorgon kan det inte regna", " "),
		strings.Split("imorgon kan det inte snöa", " "),
		strings.Split("imorgon kan det  regna", " "),
		strings.Split("imorgon kan det  snöa", " "),
		strings.Split("imorgon kommer det inte regna", " "),
		strings.Split("imorgon kommer det inte snöa", " "),
		strings.Split("imorgon kommer det  regna", " "),
		strings.Split("imorgon kommer det  snöa", " "),
	}
	if !reflect.DeepEqual(expect, result) {
		t.Errorf(fsExpGot, expect, result)
	}
}
