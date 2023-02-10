package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type line []string

// Reader struct
type Reader struct {
	inner          *csv.Reader
	CaseSensHeader bool

	allowMissingFields bool
	allowUnknownFields bool
	allowOrderMismatch bool
	requiredFields     map[string]bool

	inputHeaderSize        int
	headerStructableFields map[string]int // used for non-strict mode
}

func (r *Reader) strict() bool {
	return r.allowUnknownFields == false &&
		r.allowMissingFields == false &&
		r.allowOrderMismatch == false
}

// Strict disallows unknown fields, missing fields, and field order mismatch
func (r *Reader) Strict() {
	r.allowUnknownFields = false
	r.allowMissingFields = false
	r.allowOrderMismatch = false
}

// NonStrict allows unknown fields, missing fields, and field order mismatch
func (r *Reader) NonStrict() {
	r.allowUnknownFields = true
	r.allowMissingFields = true
	r.allowOrderMismatch = true
}

// AllowOrderMismatch allows field order to mismatch between struct and header
func (r *Reader) AllowOrderMismatch() {
	r.allowOrderMismatch = true
}

// AllowUnknownFields allows unknown fields in header
func (r *Reader) AllowUnknownFields() {
	r.allowUnknownFields = true
}

// AllowMissingFields allows missing fields in header
func (r *Reader) AllowMissingFields() {
	r.allowMissingFields = true
}

// if set, the parser will accept input lines with fewer columns than the header (is the last column is empty, some converters will skip it, hence this method could be useful)
func (r *Reader) AcceptShortLines() {
	r.inner.FieldsPerRecord = -1
}

// if set, the reader accepts headers missing any fields except for these
func (r *Reader) RequiredFields(fields ...string) {
	m := make(map[string]bool)
	for _, s := range fields {
		if !r.CaseSensHeader {
			s = strings.ToLower(s)
		}
		m[s] = true
	}
	r.requiredFields = m
}

// ReadLine reads the next line from the input data
// Returns bool, error
// - bool is true if a line was read; false if we were at the end of the file
func (r *Reader) ReadLine(v any) (bool, error) {
	fs, err := r.inner.Read()
	if err != nil {
		if err != io.EOF {
			return false, nil
		}
		return false, fmt.Errorf("failed to read line %v : %v", strings.Join(fs, string(r.inner.Comma)), err)
	}
	err = r.Unmarshal(fs, v)
	return true, err
}

func (r *Reader) ReadHeader(v any) error {
	header, err := r.inner.Read()
	if err != nil {
		return err
	}
	return r.validateHeader(header, v)
}

func (r *Reader) validateHeader(header line, v any) error {
	r.inputHeaderSize = len(header)
	s := reflect.ValueOf(v).Elem()
	r.headerStructableFields = make(map[string]int)
	if r.strict() {
		if s.NumField() != len(header) {
			return &fieldMismatch{s.NumField(), len(header)}
		}
		for i := 0; i < s.NumField(); i++ {
			ss := s.Type().Field(i).Name
			hs := header[i]
			if !r.CaseSensHeader {
				hs = strings.ToLower(hs)
				ss = strings.ToLower(ss)
			}
			if ss != hs {
				return fmt.Errorf("struct field does not match header field. Expected %v found %v", hs, ss)
			}
			r.headerStructableFields[hs] = i
		}
	} else {
		missingReqFields := []string{}
		missingFields := []string{}
		headerFields := make(map[string]int)
		for i, s := range header {
			if !r.CaseSensHeader {
				s = strings.ToLower(s)
			}
			headerFields[s] = i
		}
		structFields := map[string]int{}
		for i := 0; i < s.NumField(); i++ {
			ss := s.Type().Field(i).Name
			if !r.CaseSensHeader {
				ss = strings.ToLower(ss)
			}
			structFields[ss] = i
			hIndex, inHeader := headerFields[ss]
			if inHeader {
				r.headerStructableFields[ss] = hIndex
			} else {
				if len(r.requiredFields) > 0 && r.requiredFields[ss] {
					missingFields = append(missingFields, ss)
				} else if !r.allowMissingFields {
					missingFields = append(missingFields, ss)
				}
			}
		}
		if !r.allowUnknownFields {
			unknownFields := []string{}
			for f := range headerFields {
				if _, inStruct := structFields[f]; !inStruct {
					unknownFields = append(unknownFields, f)
				}
			}
			if len(unknownFields) > 0 {
				return fmt.Errorf("header contains unknown fields %s; found: %s", strings.Join(unknownFields, " "), strings.Join(header, " "))
			}
		}
		if !r.allowOrderMismatch {
			hfs := []string{}
			sfs := []string{}
			for _, s := range sortKeysByValue(headerFields) {
				if _, inStruct := structFields[s]; inStruct {
					hfs = append(hfs, s)
				}
			}
			for _, s := range sortKeysByValue(structFields) {
				if _, inHeader := headerFields[s]; inHeader {
					sfs = append(sfs, s)
				}
			}
			if !reflect.DeepEqual(sfs, hfs) {
				return fmt.Errorf("header is not ordered according to struct; struct: %s, header: %s", strings.Join(sfs, " "), strings.Join(hfs, " "))
			}
		}
		if len(missingReqFields) > 0 {
			return fmt.Errorf("header missing required fields %s; found: %s", strings.Join(missingFields, " "), strings.Join(header, " "))
		}
		if len(missingFields) > 0 {
			return fmt.Errorf("header missing fields %s; found: %s", strings.Join(missingFields, " "), strings.Join(header, " "))
		}
	}
	for required := range r.requiredFields {
		if _, structable := r.headerStructableFields[required]; !structable {
			return fmt.Errorf("required field %s does not exist in struct", required)
		}
	}
	return nil
}

func NewReader(source io.Reader, separator rune) *Reader {
	r := Reader{inner: csv.NewReader(source)}
	r.inner.Comma = separator
	r.inner.LazyQuotes = true
	return &r
}

func NewStringReader(source string, separator rune) *Reader {
	return NewReader(strings.NewReader(source), separator)
}

func NewFileReader(fName string, separator rune) (*Reader, error) {
	file, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	return NewReader(file, separator), nil
}

func (r *Reader) Unmarshal(line []string, v any) error {
	struc := reflect.ValueOf(v).Elem()
	if r.inner.FieldsPerRecord > 0 && r.inputHeaderSize != len(line) {
		return &fieldMismatch{r.inputHeaderSize, len(line)}
	}
	if r.inner.FieldsPerRecord < 1 {
		for len(line) < struc.NumField() {
			line = append(line, "")
		}
	}
	for i := 0; i < struc.NumField(); i++ {
		f := struc.Field(i)
		name := struc.Type().Field(i).Name
		if !r.CaseSensHeader {
			name = strings.ToLower(name)
		}
		colIndex, structableFields := r.headerStructableFields[name]
		if !structableFields {
			continue
		}
		val := line[colIndex]
		switch f.Type().String() {
		case "string":
			f.SetString(val)
		case "int":
			if val == "" {
				return fmt.Errorf("empty int field %s for input line %v", name, strings.Join(line, string(r.inner.Comma)))
			}
			ival, err := strconv.ParseInt(val, 10, 0)
			if err != nil {
				return err
			}
			f.SetInt(ival)
		case "bool":
			if val == "" {
				return fmt.Errorf("empty boolean field %s for input line %v", name, strings.Join(line, string(r.inner.Comma)))
			}
			bval, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			f.SetBool(bval)
		default:
			return fmt.Errorf("unsupported type: " + f.Type().String())
		}
	}
	return nil
}

func sortKeysByValue(m map[string]int) []string {
	values := maps.Values(m)
	slices.Sort(values)
	res := []string{}
	for _, v := range values {
		for k, v0 := range m {
			if v == v0 {
				res = append(res, k)
			}
		}
	}
	return res
}

// re-usable error types
type fieldMismatch struct {
	expected, found int
}

func (e *fieldMismatch) Error() string {
	return "csv line fields mismatch; expected " + strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}
