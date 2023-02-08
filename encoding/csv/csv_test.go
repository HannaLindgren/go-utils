package csv

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"
)

var staticImport = fmt.Sprintf("import fmt") // always keep fmt in import list

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
	reader.Strict = true
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
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
			return
		}
		res = append(res, entry)
	}
	if len(res) != 2 {
		t.Errorf("Expected %v entries, got %v: %#v", 2, len(res), res)
		return
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
		return
	}
	reader.Strict = true
	var header entry
	err = reader.ReadHeader(&header)
	if err == nil {
		t.Errorf("Expected error here")
	}
}

func TestCsvReaderNonStrictCaseSens(t *testing.T) {
	var err error
	var source = `Country	OrigLang	Orth	Exonym	Priority	Checked	Comment	Template
GBR	eng	The Thames	Themsen	4	true	hepp	auto_case
BEL	fre	Bruxelles	Bryssel	3	false		hardwired_cities`
	var separator = '	'
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = true
	if err != nil {
		t.Errorf("Got error from NewStringReader: %v", err)
		return
	}
	reader.Strict = false
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	for {
		var entry entry
		err := reader.Read(&entry)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		// bts, err := json.Marshal(entry)
		// if err != nil {
		// 	t.Errorf("Got error from json.Marshal: %v", err)
		// }
		//fmt.Println(string(bts))
	}
}

func TestCsvReaderNonStrictCaseInsens(t *testing.T) {
	var err error
	var source = `country	OrigLang	Orth	Exonym	Priority	Checked	Comment	Template
GBR	eng	The Thames	Themsen	4	true	hepp	auto_case
BEL	fre	Bruxelles	Bryssel	3	false		hardwired_cities`
	var separator = '	'
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = false
	reader.Strict = false
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	for {
		var entry entry
		err := reader.Read(&entry)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		// bts, err := json.Marshal(entry)
		// if err != nil {
		// 	t.Errorf("Got error from json.Marshal: %v", err)
		// }
		//fmt.Println(string(bts))
	}
}

func TestCsvReaderNonStrictCaseInsensShortLines(t *testing.T) {
	var err error
	var source = `country	OrigLang	Orth	Exonym	Priority	Checked	Comment	Template
GBR	eng	The Thames	Themsen	4	true	hepp	auto_case
BEL	fre	Bruxelles	Bryssel	3	false		hardwired_cities
FRA	fre	Paris		3	false`
	var separator = '	'
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = false
	reader.Strict = false
	reader.AcceptShortLines()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	for {
		var entry entry
		err := reader.Read(&entry)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		// bts, err := json.Marshal(entry)
		// if err != nil {
		// 	t.Errorf("Got error from json.Marshal: %v", err)
		// }
		//fmt.Println(string(bts))
	}
}

func TestCsvReaderWithMissingFields(t *testing.T) {
	var err error
	var source = `Country	OrigLang	Orth	Exonym	Priority	Checked
GBR	eng	The Thames	Themsen	4	true
BEL	fre	Bruxelles	Bryssel	3	false
FRA	fre	Paris		3	false`
	var separator = '	'
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = true
	reader.Strict = false
	reader.AcceptShortLines()
	reader.AcceptMissingFields("Comment")
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	for {
		var entry entry
		err := reader.Read(&entry)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		// bts, err := json.Marshal(entry)
		// if err != nil {
		// 	t.Errorf("Got error from json.Marshal: %v", err)
		// }
		//fmt.Println(string(bts))
	}
}

func TestCsvReaderWithRequiredFields(t *testing.T) {
	var err error
	var source = `Country	OrigLang	Orth	Exonym	Priority	Checked
GBR	eng	The Thames	Themsen	4	true
BEL	fre	Bruxelles	Bryssel	3	false
FRA	fre	Paris		3	false`
	var separator = '	'
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = true
	reader.Strict = false
	reader.AcceptShortLines()
	reader.RequiredFields("Orig", "OrigLang")
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	for {
		var entry entry
		err := reader.Read(&entry)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		// bts, err := json.Marshal(entry)
		// if err != nil {
		// 	t.Errorf("Got error from json.Marshal: %v", err)
		// }
		//fmt.Println(string(bts))
	}
}
