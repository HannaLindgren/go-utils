package util

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
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

// ConvertLinesFromStinOrFiles takes a command name (for logging), an array of arguments (typically from os.Args), and a conversion function, that converts an input string to another (output string). Utility for writing simple code for processing textfiles, typically converting each line into another output line (upcase, line length, etc).
func ConvertLinesFromStinOrFiles(cmdname string, args []string, convert fn) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <files or strings>\n", cmdname)
		os.Exit(1)
	}
	for i := 1; i < len(args); i++ {
		a := args[i]
		if _, err := os.Stat(a); os.IsNotExist(err) {
			fmt.Println(convert(a))
		} else {
			r, fh, err := GetFileReader(a)
			defer fh.Close()
			if err != nil {
				log.Fatalf("%v", err)
			}
			scan := bufio.NewScanner(r)
			for scan.Scan() {
				s := scan.Text()
				fmt.Println(convert(s))
			}
		}
	}

}
