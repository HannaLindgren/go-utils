package util

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// GetFileReader reads an input file, gzipped or plain text, and returns an io.Reader for line scanning, along with the file handle, that needs to be closed after reading.
func GetFileReader(fName string) (io.Reader, *os.File, error) {
	fh, err := os.Open(fName)
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

// ConvertAndPrintFromFileArgs takes a conversion function and an array of files to convert. The conversion function should convert an input string to another (output) string. It's a utility for writing simple code for processing textfiles, typically converting each input line into another output line (upcase, line length, etc).
func ConvertAndPrintFromFileArgs(convert fn, args []string) error {
	for i := 0; i < len(args); i++ {
		f := args[i]
		r, fh, err := GetFileReader(f)
		defer fh.Close()
		if err != nil {
			return err
		}
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			s := scan.Text()
			fmt.Println(convert(s))
		}
	}
	return nil
}

// ConvertAndPrintFromFileArgsOrStdin takes a conversion function, and as conversion input it uses (1) files specified in os.Args; or (2) stdin. The conversion function should convert an input string to another (output) string. It's a utility for writing simple code for processing textfiles, typically converting each input line into another output line (upcase, line length, etc).
func ConvertAndPrintFromFileArgsOrStdin(convert fn) error {
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			f := os.Args[i]
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
