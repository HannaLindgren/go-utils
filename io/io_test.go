package io

import (
	"testing"
)

var fsExpGot = "expected: %#v ; got: %#v"

func TestDetachFileExtension(t *testing.T) {
	var fName, expectFN, expectExt, resultFN, resultExt string

	fName = "/tmp/afilename.txt"
	expectFN = "/tmp/afilename"
	expectExt = "txt"
	resultFN, resultExt = DetachFileExtension(fName)
	if resultFN != expectFN {
		t.Errorf(fsExpGot, expectFN, resultFN)
	}
	if resultExt != expectExt {
		t.Errorf(fsExpGot, expectExt, resultExt)
	}
}
