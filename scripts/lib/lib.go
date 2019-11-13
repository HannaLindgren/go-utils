package lib

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// IsFile returns true if the given file exists (as a file or as a directory)
func IsFile(fName string) bool {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetFileReader reads an input file, gzipped or plain text, and returns an io.Reader for line scanning, along with the file handle, that needs to be closed after reading.
func GetFileReader(fName string) (io.Reader, *os.File, error) {
	fh, err := os.Open(filepath.Clean(fName))
	//defer fh.Close()
	if err != nil {
		return nil, fh, fmt.Errorf("Couldn't open file %s for reading : %v\n", fName, err)
	}

	if strings.HasSuffix(fName, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			return nil, fh, fmt.Errorf("Couldn't to open gz reader : %v", err)
		}
		return io.Reader(gz), fh, nil
	}
	return io.Reader(fh), fh, nil
}

type fn func(string) string

// func ConvertAndPrintStdin(convert fn) error {
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		s := scanner.Text()
// 		fmt.Println(convert(s))
// 	}
// 	return nil
// }

// ReadFileToString
func ReadFileToString(fName string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Clean(fName))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadStdinToString
func ReadStdinToString() (string, error) {
	stdin := bufio.NewReader(os.Stdin)
	b, err := ioutil.ReadAll(stdin)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ConvertAndPrintFromFilesOrStdin
func ConvertAndPrintFromFilesOrStdin(convert fn, files []string) error {
	if len(files) > 0 {
		for _, f := range files {
			r, fh, err := GetFileReader(f)
			defer fh.Close()
			if err != nil {
				return err
			}
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				s := scanner.Text()
				fmt.Println(convert(s))
			}

		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			s := scanner.Text()
			fmt.Println(convert(s))
		}
	}
	return nil
}

// ConvertAndPrintFromFileArgsOrStdin takes a conversion function, and as conversion input it uses (1) files specified in os.Args; or (2) stdin. The conversion function should convert an input string to another (output) string. It's a utility for writing simple code for processing textfiles, typically converting each input line into another output line (upcase, line length, etc).
func ConvertAndPrintFromFileArgsOrStdin(convert fn) error {
	args := os.Args[1:]
	if len(args) > 0 {
		for _, arg := range args {
			if IsFile(arg) {
				r, fh, err := GetFileReader(arg)
				defer fh.Close()
				if err != nil {
					return err
				}
				scanner := bufio.NewScanner(r)
				for scanner.Scan() {
					s := scanner.Text()
					fmt.Println(convert(s))
				}
			} else {
				fmt.Println(convert(arg))
			}

		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			s := scanner.Text()
			fmt.Println(convert(s))
		}
	}
	return nil
	//return ConvertAndPrintFromFilesOrStdin(convert, os.Args[1:])
}
