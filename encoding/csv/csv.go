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

	// Strict true: header and struct must match exactly
	// Strict false: struct can be a subset of the header defined columns [default]
	Strict    bool
	headerDef map[int]string // used for non-strict mode
}

// todo: better handling of io.EOF error
func (r *Reader) Read(v interface{}) error {
	fs, err := r.inner.Read()
	if err != nil {
		return err
	}
	return r.Unmarshal(fs, v)
	return nil
}

func (r *Reader) ReadHeader(v interface{}) error {
	header, err := r.inner.Read()
	if err != nil {
		return err
	}
	r.headerDef = make(map[int]string)
	for i, s := range header {
		r.headerDef[i] = s
	}
	s := reflect.ValueOf(v).Elem()

	if r.Strict {
		if s.NumField() != len(header) {
			return &FieldMismatch{s.NumField(), len(header)}
		}
		for i := 0; i < s.NumField(); i++ {
			s := s.Type().Field(i).Name
			h := header[i]
			if r.CaseSensHeader {
				if s != h {
					return &FieldNameMismatch{s, h}
				}
			} else if !strings.EqualFold(s, h) {
				return &FieldNameMismatch{s, h}
			}
		}
	} else {
		return fmt.Errorf("Non-strict mode is not implemented")
	}
	return nil
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

func (r *Reader) Unmarshal(line []string, v interface{}) error {
	s := reflect.ValueOf(v).Elem()
	if r.Strict {
		if s.NumField() != len(line) {
			return &FieldMismatch{s.NumField(), len(line)}
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
	} else {
		return fmt.Errorf("Non-strict mode is not implemented")
	}
	return nil
}
