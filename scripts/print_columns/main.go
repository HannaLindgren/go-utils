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

func process(requestedFields map[string]int, lines []string) error {
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
			colsToPrint[i] = ri
			nColsToPrint++
		}
	}

	// check for invalid columns in input flag
	for f := range requestedFields {
		if !slices.Contains(existingFields, f) {
			return fmt.Errorf("requested field %s does not exist in input", f)
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

func main() {

	caseSens = flag.Bool("case", false, "Case sensitive column headers")
	fieldSepFlag := flag.String("sep", "<tab>", "Field `separator`")
	skipHeader = flag.Bool("skip", false, "Do not include header in output")
	fieldSep = *fieldSepFlag
	if fieldSep == "<tab>" {
		fieldSep = "\t"
	}

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Print selected columns based on file header")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <requested columns> <files>*\n", cmdname)
		fmt.Fprintf(os.Stderr, "       OR\n")
		fmt.Fprintf(os.Stderr, "cat <files> | %s <requested columns>\n", cmdname)
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Example usage:\n")
		fmt.Fprintf(os.Stderr, "%s orig,NormName /tmp/prio4.txt\n", cmdname)
	}

	flag.Parse()

	if flag.NArg() < 2 {
		printUsage()
		os.Exit(0)
	}

	var requestedFields = map[string]int{}
	for i, f := range columnSplitRE.Split(flag.Args()[0], -1) {
		if !*caseSens {
			f = strings.ToLower(f)
		}
		requestedFields[f] = i
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
