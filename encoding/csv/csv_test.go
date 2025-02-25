package csv

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

var staticImport = fmt.Sprintf("import fmt") // always keep fmt in import list

// structs for testing
type entry struct {
	Country  string
	OrigLang string
	Orth     string
	Exonym   string
	Priority int
	Checked  bool
	Comment  string
}

type entryWithChildren struct {
	Country  string
	OrigLang string
	Orth     string
	Exonym   string
	Priority int
	Checked  bool
	Comment  string
	Children []string
}

func TestCsvReaderUseCase(t *testing.T) {
	var err error
	var source = `country	origLang	orth	exonym	priority	checked	comment
GBR	eng	The Thames	Themsen	4	true	hepp
BEL	fre	Bruxelles	Bryssel"	3	false	`
	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.Strict()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	res := []entry{}
	for {
		var entry entry
		hasNext, err := reader.ReadLine(&entry)
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
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

func TestCsvReaderUseCaseWithChildren(t *testing.T) {
	var err error
	var source = `country	origLang	orth	exonym	priority	checked	comment
GBR	eng	The Thames	Themsen	4	true	hepp
SWE	swe	Mälaren	Lake Mälaren	1	true	
BEL	fre	Bruxelles	Bryssel"	3	false	`
	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.AllowMissingFields()
	var header entryWithChildren
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	nEntries := 0
	for {
		var entry entryWithChildren
		hasNext, err := reader.ReadLine(&entry)
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
		}
		nEntries++
	}
	if nEntries != 3 {
		t.Errorf("Expected %v entries, got %v", 3, nEntries)
		return
	}
}

func TestCsvReaderWithInvalidHeader(t *testing.T) {
	var err error
	var source = `country	origLang	orth	exonym	priority	checked	comments
GBR	eng	The Thames	Themsen	4	true	hepp
BEL	fre	Bruxelles	Bryssel	3	false	`
	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.Strict()
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
	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = true
	if err != nil {
		t.Errorf("Got error from NewStringReader: %v", err)
		return
	}
	reader.Strict()
	var header entry
	err = reader.ReadHeader(&header)
	if err == nil {
		t.Errorf("Expected error here")
	}
}

func TestCsvReaderNonStrictCaseSens(t *testing.T) {
	var err error
	var source = `Country	OrigLang	Orth	Exonym	Template	Priority	Checked	Comment
GBR	eng	The Thames	Themsen	tmpl	4	true	hepp
BEL	fre	Bruxelles	Bryssel	tmpl	3	false	`
	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = true
	reader.AllowUnknownFields()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	for {
		var entry entry
		hasNext, err := reader.ReadLine(&entry)
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
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
	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = false
	reader.AllowUnknownFields()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	for {
		var entry entry
		hasNext, err := reader.ReadLine(&entry)
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
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
	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = false
	reader.AllowUnknownFields()
	reader.AcceptShortLines()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	for {
		var entry entry
		hasNext, err := reader.ReadLine(&entry)
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
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
	var source = `Country	Template	OrigLang	Orth	Exonym	Checked
GBR	tmpl2	eng	The Thames	Themsen	true
BEL	tmpl1	fre	Bruxelles	Bryssel	false
	tmpl3	fre	Paris		false`
	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.CaseSensHeader = true
	reader.AllowMissingFields()
	reader.AllowUnknownFields()
	//reader.AcceptShortLines()
	reader.RequiredFields("Orth", "OrigLang")
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		t.Errorf("Got error from ReadHeader: %v", err)
		return
	}
	for {
		var entry entry
		hasNext, err := reader.ReadLine(&entry)
		if err != nil {
			t.Errorf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
		}
		bts, err := json.Marshal(entry)
		if err != nil {
			t.Errorf("Got error from json.Marshal: %v", err)
		}
		fmt.Println(string(bts))
	}
}

func ExampleReader_ReadLine_strict() {
	var err error
	var source = `country	origLang	orth	exonym	priority	checked	comment
GBR	eng	The Thames	Themsen	4	true	todo
BEL	fre	Bruxelles	Bryssel"	3	false	`

	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.Strict()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		log.Printf("Got error from ReadHeader: %v", err)
		return
	}
	res := []entry{}
	for {
		var e entry
		hasNext, err := reader.ReadLine(&e)
		if err != nil {
			log.Printf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
		}
		res = append(res, e)
	}
	for _, e := range res {
		bts, err := json.Marshal(e)
		if err != nil {
			log.Printf("Got error from json.Marshal: %v", err)
		}
		fmt.Println(string(bts))
	}

	// Output: {"Country":"GBR","OrigLang":"eng","Orth":"The Thames","Exonym":"Themsen","Priority":4,"Checked":true,"Comment":"todo"}
	// {"Country":"BEL","OrigLang":"fre","Orth":"Bruxelles","Exonym":"Bryssel\"","Priority":3,"Checked":false,"Comment":""}
}

func ExampleReader_ReadLine_allowOrderMismatch() {
	var err error
	var source = `origLang	country	orth	exonym	priority	checked	comment
GBR	eng	The Thames	Themsen	4	true	todo
BEL	fre	Bruxelles	B"ryssel"	3	false	`

	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.AllowOrderMismatch()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		log.Printf("Got error from ReadHeader: %v", err)
		return
	}
	res := []entry{}
	for {
		var e entry
		hasNext, err := reader.ReadLine(&e)
		if err != nil {
			log.Printf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
		}
		res = append(res, e)
	}
	for _, e := range res {
		bts, err := json.Marshal(e)
		if err != nil {
			log.Printf("Got error from json.Marshal: %v", err)
		}
		fmt.Println(string(bts))
	}

	// Output: {"Country":"eng","OrigLang":"GBR","Orth":"The Thames","Exonym":"Themsen","Priority":4,"Checked":true,"Comment":"todo"}
	// {"Country":"fre","OrigLang":"BEL","Orth":"Bruxelles","Exonym":"B\"ryssel\"","Priority":3,"Checked":false,"Comment":""}
}

func ExampleReader_ReadLine_nonStrict() {
	var err error
	var source = `origLang	country	orth	exonym	priority	template	comment
GBR	eng	The Thames	Themsen	4	tmpl1	todo
BEL	fre	Bruxelles	Bryssel	3	tmpl2	`

	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.NonStrict()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		log.Printf("Got error from ReadHeader: %v", err)
		return
	}
	res := []entry{}
	for {
		var e entry
		hasNext, err := reader.ReadLine(&e)
		if err != nil {
			log.Printf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
		}
		res = append(res, e)
	}
	for _, e := range res {
		bts, err := json.Marshal(e)
		if err != nil {
			log.Printf("Got error from json.Marshal: %v", err)
		}
		fmt.Println(string(bts))
	}

	// Output: {"Country":"eng","OrigLang":"GBR","Orth":"The Thames","Exonym":"Themsen","Priority":4,"Checked":false,"Comment":"todo"}
	// {"Country":"fre","OrigLang":"BEL","Orth":"Bruxelles","Exonym":"Bryssel","Priority":3,"Checked":false,"Comment":""}
}

func TestReadLine_strictButAllowOrderMismatch(t *testing.T) {
	var err error
	var source = `origLang	country	orth	exonym	priority	checked	template	comment
GBR	eng	The Thames	Themsen	4	true	tmpl1	todo
BEL	fre	Bruxelles	Bryssel"	3	false	tmpl2	`

	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.Strict()
	reader.AllowOrderMismatch()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		expectSubstring := "header contains unknown field"
		if !strings.Contains(fmt.Sprintf("%v", err), expectSubstring) {
			t.Errorf("expected error: %v, found: %v", expectSubstring, err)

		}
	} else {
		t.Errorf("expected error here")
	}
}

func TestReadLine_nonStrictButDisallowOrderMismatch(t *testing.T) {
	var err error
	var source = `origLang	country	orth	exonym	priority	checked	template	comment
GBR	eng	The Thames	Themsen	4	true	tmpl1	todo
BEL	fre	Bruxelles	Bryssel"	3	false	tmpl2	`

	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.Strict()
	reader.AllowUnknownFields()
	var header entry
	err = reader.ReadHeader(&header)
	if err != nil {
		expectSubstring := "header is not ordered according to struct"
		if !strings.Contains(fmt.Sprintf("%v", err), expectSubstring) {
			t.Errorf("expected error: %v, found: %v", expectSubstring, err)

		}
	} else {
		t.Errorf("expected error here")
	}
}

type entryWithTagsCaseSens struct {
	Country  string `csv:"country"`
	OrigLang string `csv:"origLang"`
	Orth     string `csv:"orth"`
	Exonym   string `csv:"exonym"`
	Priority int    `csv:"Priority"`
	Checked  bool   `csv:"checked"`
	Comment  string `csv:"translator comment"`
}

func ExampleReader_ReadLine_strict_Tags() {
	var err error
	var source = `country	origLang	orth	exonym	Priority	checked	translator comment
GBR	eng	The Thames	Themsen	4	true	todo
BEL	fre	Bruxelles	Brysseles	3	false	`

	var separator = "\t"
	var reader = NewStringReader(source, separator)
	reader.Strict()
	reader.CaseSensHeader = true
	var header entryWithTagsCaseSens
	err = reader.ReadHeader(&header)
	if err != nil {
		log.Printf("Got error from ReadHeader: %v", err)
		return
	}
	res := []entryWithTagsCaseSens{}
	for {
		var e entryWithTagsCaseSens
		hasNext, err := reader.ReadLine(&e)
		if err != nil {
			log.Printf("Got error from Read: %v", err)
			return
		}
		if !hasNext {
			break
		}
		res = append(res, e)
	}
	for _, e := range res {
		bts, err := json.Marshal(e)
		if err != nil {
			log.Printf("Got error from json.Marshal: %v", err)
		}
		fmt.Println(string(bts))
	}

	// Output: {"Country":"GBR","OrigLang":"eng","Orth":"The Thames","Exonym":"Themsen","Priority":4,"Checked":true,"Comment":"todo"}
	// {"Country":"BEL","OrigLang":"fre","Orth":"Bruxelles","Exonym":"Brysseles","Priority":3,"Checked":false,"Comment":""}
}

// type entryWithTagsCaseInsens struct {
// 	Country  string `csv:"country"`
// 	OrigLang string `csv:"origLang"`
// 	Orth     string `csv:"orth"`
// 	Exonym   string `csv:"exonym"`
// 	Priority int    `csv:"Priority"`
// 	Checked  bool   `csv:"checked"`
// 	Comment  string `csv:"translator Comment"`
// }

// func ExampleReader_ReadLine_NonStrict_Tags() {
// 	var err error
// 	var source = `country	OrigLang	orth	exonym	Priority	checked	translator comment
// GBR	eng	The Thames	Themsen	4	true	todox
// BEL	fre	Bruxelles	Brysseles	3	false	`

// 	var separator = "\t"
// 	var reader = NewStringReader(source, separator)
// 	reader.NonStrict()
// 	reader.CaseSensHeader = false
// 	var header entryWithTagsCaseInsens
// 	err = reader.ReadHeader(&header)
// 	if err != nil {
// 		log.Printf("Got error from ReadHeader: %v", err)
// 		return
// 	}
// 	res := []entryWithTagsCaseInsens{}
// 	for {
// 		var e entryWithTagsCaseInsens
// 		hasNext, err := reader.ReadLine(&e)
// 		if err != nil {
// 			log.Printf("Got error from Read: %v", err)
// 			return
// 		}
// 		if !hasNext {
// 			break
// 		}
// 		res = append(res, e)
// 	}
// 	for _, e := range res {
// 		bts, err := json.Marshal(e)
// 		if err != nil {
// 			log.Printf("Got error from json.Marshal: %v", err)
// 		}
// 		fmt.Println(string(bts))
// 	}

// 	// Output: {"Country":"GBR","OrigLang":"eng","Orth":"The Thames","Exonym":"Themsen","Priority":4,"Checked":true,"Comment":"todox"}
// 	// {"Country":"BEL","OrigLang":"fre","Orth":"Bruxelles","Exonym":"Brysseles","Priority":3,"Checked":false,"Comment":""}
// }
