package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"

	hio "github.com/HannaLindgren/go-utils/io"
)

func process(requestedFields map[string]field, lines []string) error {
	if len(lines) == 0 {
		return fmt.Errorf("no input lines")
	}
	header := lines[0]
	if !*caseSens {
		header = strings.ToLower(header)
	}
	existingFields := strings.Split(header, fieldSep)

	// save printable indices
	var colsToPrint = make(map[int]int)
	var nColsToPrint = 0
	for i, s := range existingFields {
		if ri, doPrint := requestedFields[s]; doPrint {
			colsToPrint[i] = ri.index
			nColsToPrint++
		}
	}

	// check for invalid columns in input flag
	for key, field := range requestedFields {
		if !slices.Contains(existingFields, key) {
			return fmt.Errorf("requested field %s does not exist in input data", field.name)
		}
	}

	for li, l := range lines {
		if li == 0 && *skipHeader {
			continue
		}
		outFS := []string{}
		for len(outFS) < nColsToPrint {
			outFS = append(outFS, "")
		}
		for i, f := range strings.Split(l, fieldSep) {
			if outI, doPrint := colsToPrint[i]; doPrint {
				outFS[outI] = f
			}
		}
		fmt.Println(strings.Join(outFS, "\t"))
	}
	return nil
}

const cmdname = "print_columns"

// options
var caseSens, skipHeader *bool
var fieldSep string

var columnSplitRE = regexp.MustCompile("[,;: ]+")

type field struct {
	index int
	name  string
}

func main() {

	caseSens = flag.Bool("c", false, "Case sensitive column headers")
	fieldSepFlag := flag.String("s", "<tab>", "Field `separator`")
	skipHeader = flag.Bool("H", false, "Skip output header")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Print selected columns based on file header")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <requested columns> <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       OR\n")
		fmt.Fprintf(os.Stderr, "cat <files> | %s <requested columns>\n", cmdname)
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Example usage:\n")
		fmt.Fprintf(os.Stderr, "%s orth,country /tmp/sourcefile.txt\n", cmdname)
	}

	flag.Parse()

	if flag.NArg() != 1 && flag.NArg() != 2 {
		printUsage()
		os.Exit(0)
	}

	fieldSep = *fieldSepFlag
	if fieldSep == "<tab>" {
		fieldSep = "\t"
	}

	var requestedFields = map[string]field{}
	for i, f := range columnSplitRE.Split(flag.Args()[0], -1) {
		var key = f
		if !*caseSens {
			key = strings.ToLower(f)
		}
		requestedFields[key] = field{index: i, name: f}
	}

	//fmt.Fprintf(os.Stderr, "%#v\n", requestedFields)

	if flag.NArg() > 1 {
		for _, f := range flag.Args()[1:] {
			lines, err := hio.ReadFileToLines(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read from file %s: %v\n", f, err)
				os.Exit(1)
			}
			err = process(requestedFields, lines)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to process file %s: %v\n", f, err)
				os.Exit(1)
			}
		}
	} else { // read from stdin
		lines, err := hio.ReadStdinToLines()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read from stdin: %v\n", err)
			os.Exit(1)
		}
		err = process(requestedFields, lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to process stdin: %v\n", err)
			os.Exit(1)
		}
	}

}
