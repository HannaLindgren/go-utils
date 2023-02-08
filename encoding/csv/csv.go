package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// error types
type FieldMismatch struct {
	expected, found int
}

func (e *FieldMismatch) Error() string {
	return "CSV line fields mismatch. Expected " + strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}

type FieldNameMismatch struct {
	expected, found string
}

func (e *FieldNameMismatch) Error() string {
	return fmt.Sprintf("Struct field does not match header field. Expected %v found %v", e.expected, e.found)
}

type UnsupportedType struct {
	Type string
}

func (e *UnsupportedType) Error() string {
	return "Unsupported type: " + e.Type
}

type line []string

// Reader struct
type Reader struct {
	inner          *csv.Reader
	CaseSensHeader bool

	acceptMissingFields map[string]bool

	// Strict true: header and struct must match exactly
	// Strict false: struct can be a subset of the header defined columns [default]
	Strict    bool
	headerDef map[int]bool // used for non-strict mode
}

// todo: better handling of io.EOF error
func (r *Reader) Read(v interface{}) error {
	fs, err := r.inner.Read()
	if err != nil {
		if err != io.EOF {
			return fmt.Errorf("failed to read line %v : %v", strings.Join(fs, string(r.inner.Comma)), err)
		}
		return err
	}
	return r.Unmarshal(fs, v)
}

func (r *Reader) ReadHeader(v interface{}) error {
	header, err := r.inner.Read()
	if err != nil {
		return err
	}
	return r.validateHeader(header, v)
}

func (r *Reader) validateHeader(header line, v interface{}) error {
	s := reflect.ValueOf(v).Elem()
	r.headerDef = make(map[int]bool)
	if r.Strict {
		if s.NumField() != len(header) {
			return &FieldMismatch{s.NumField(), len(header)}
		}
		for i := 0; i < s.NumField(); i++ {
			ss := s.Type().Field(i).Name
			hs := header[i]
			if r.CaseSensHeader {
				if ss != hs {
					return &FieldNameMismatch{ss, hs}
				}
			} else if !strings.EqualFold(ss, hs) {
				return &FieldNameMismatch{ss, hs}
			}
			r.headerDef[i] = true
		}
	} else {
		missingFields := []string{}
		hMap := make(map[string]int)
		for i, s := range header {
			if !r.CaseSensHeader {
				s = strings.ToLower(s)
			}
			hMap[s] = i
		}
		for i := 0; i < s.NumField(); i++ {
			ss := s.Type().Field(i).Name
			if !r.CaseSensHeader {
				ss = strings.ToLower(ss)
			}
			hIndex, inHeader := hMap[ss]
			if inHeader {
				r.headerDef[hIndex] = true
			} else {
				if !r.acceptMissingFields[ss] {
					missingFields = append(missingFields, ss)
				}
			}
		}
		if len(missingFields) > 0 {
			return fmt.Errorf("Header missing struct fields %s; found: %s", strings.Join(missingFields, " "), strings.Join(header, " "))
		}
	}
	return nil
}

// if set, the parser will accept input lines with fewer columns than earlier lines (is the last column is empty, some converters will skip it, hence this method could be useful)
func (r *Reader) AcceptShortLines() {
	r.inner.FieldsPerRecord = -1
}

// if set, the reader accepts headers missing fields contained in the struct
func (r *Reader) AcceptMissingFields(fields ...string) {
	m := make(map[string]bool)
	for _, s := range fields {
		if !r.CaseSensHeader {
			s = strings.ToLower(s)
		}
		m[s] = true
	}
	r.acceptMissingFields = m
}

func NewReader(source io.Reader, separator rune) *Reader {
	r := Reader{inner: csv.NewReader(source)}
	r.inner.Comma = separator
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

func (r *Reader) strictUnmarshal(line []string, v interface{}) error {
	s := reflect.ValueOf(v).Elem()
	if r.inner.FieldsPerRecord > 0 && s.NumField() != len(line) {
		return &FieldMismatch{s.NumField(), len(line)}
	}
	if r.inner.FieldsPerRecord < 1 {
		for len(line) < s.NumField() {
			line = append(line, "")
		}
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		switch f.Type().String() {
		case "string":
			f.SetString(line[i])
		case "int":
			ival, err := strconv.ParseInt(line[i], 10, 0)
			if err != nil {
				return err
			}
			f.SetInt(ival)
		case "bool":
			bval, err := strconv.ParseBool(line[i])
			if err != nil {
				return err
			}
			f.SetBool(bval)
		default:
			return &UnsupportedType{f.Type().String()}
		}
	}
	return nil
}

func (r *Reader) Unmarshal(line []string, v interface{}) error {
	if r.Strict {
		return r.strictUnmarshal(line, v)
	}
	strict := []string{}
	for i, s := range line {
		if r.headerDef[i] {
			strict = append(strict, s)
		}
	}
	return r.strictUnmarshal(strict, v)
}
