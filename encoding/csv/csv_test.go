package csv

import (
	"encoding/json"
	"io"
	"testing"
)

// struct for testing
type entry struct {
	Country  string
	OrigLang string
	Orth     string
	Exonym   string
	Priority int
	Checked  bool
	Comment  string
}

func TestCsvReaderUseCase(t *testing.T) {
	var err error
	var source = `country	origLang	orth	exonym	priority	checked	comment
GBR	eng	The Thames	Themsen	4	true	hepp
BEL	fre	Bruxelles	Bryssel	3	false	`
	var separator = '	'
	var reader = NewStringReader(source, separator)
	if err != nil {
		t.Errorf("Got error from NewStringReader: %v", err)
	}
	reader.Strict = true
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Header parse failed: %v", err)
	}
	res := []entry{}
	for {
		var entry entry
		err := reader.Read(&entry)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
		}
		res = append(res, entry)
	}
	if len(res) != 2 {
		t.Errorf("Expected %v entries, got %v: %#v", 2, len(res), res)
	}
	for _, e := range res {
		_, err := json.Marshal(e)
		if err != nil {
			t.Errorf("Got error from json.Marshal: %v", err)
		}
	}
}

func TestCsvReaderWithInvalidHeader(t *testing.T) {
	var err error
	var source = `country	origLang	orth	exonym	priority	checked	comments
GBR	eng	The Thames	Themsen	4	true	hepp
BEL	fre	Bruxelles	Bryssel	3	false	`
	var separator = '	'
	var reader = NewStringReader(source, separator)
	if err != nil {
		t.Errorf("Got error from NewStringReader: %v", err)
	}
	reader.Strict = true
	var header entry
	err = reader.ReadHeader(&header)
	if err == nil {
		t.Errorf("Expected error here")
	}
}

func TestCsvReaderWithInvalidCaseHeader(t *testing.T) {
	var err error
	var source = `country	origLang	orth	exonym	priority	checked	comment
GBR	eng	The Thames	Themsen	4	true	hepp
BEL	fre	Bruxelles	Bryssel	3	false	`
	var separator = '	'
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = true
	if err != nil {
		t.Errorf("Got error from NewStringReader: %v", err)
	}
	reader.Strict = true
	var header entry
	err = reader.ReadHeader(&header)
	if err == nil {
		t.Errorf("Expected error here")
	}
}

func TestCsvReaderNonStrict(t *testing.T) {
	var err error
	var source = `country	origLang	orth	exonym	priority	checked	comment	template
GBR	eng	The Thames	Themsen	4	true	hepp	auto_case
BEL	fre	Bruxelles	Bryssel	3	false	hardwired_cities`
	var separator = '	'
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = true
	if err != nil {
		t.Errorf("Got error from NewStringReader: %v", err)
	}
	reader.Strict = false
	var header entry
	err = reader.ReadHeader(&header)
	if err == nil {
		t.Errorf("Got error from ReadHeader: %v", err)
	}
	//t.Errorf("Test not implemented")
}
