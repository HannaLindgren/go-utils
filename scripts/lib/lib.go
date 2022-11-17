package lib

import (
	"bufio"
	"fmt"
	"os"

	"github.com/HannaLindgren/go-utils/io"
)

type convertFunc func(string) string

// ConvertAndPrintFromFilesOrStdin
func ConvertAndPrintFromFilesOrStdin(convert convertFunc, files []string) error {
	if len(files) > 0 {
		for _, f := range files {
			r, fh, err := io.GetFileReader(f)
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

// ConvertAndPrintFromArgsOrStdin takes a conversion function, and as conversion input it uses (1) files or strings specified in args; or (2) stdin. The conversion function should convert an input string to another (output) string. It's a utility for writing simple code for processing textfiles, typically converting each input line into another output line (upcase, line length, etc).
func ConvertAndPrintFromArgsOrStdin(convert convertFunc, args []string) error {
	if len(args) > 0 {
		for _, arg := range args {
			if io.IsFile(arg) {
				r, fh, err := io.GetFileReader(arg)
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
