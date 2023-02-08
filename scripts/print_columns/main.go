package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	hio "github.com/HannaLindgren/go-utils/io"
)

type field struct {
	req      reqField
	outIndex int
	inIndex  int
}

type reqField struct {
	index    int
	normName string
	name     string
}

var printedHeaders = map[string]bool{}

func process(requestedFields []reqField, lines []string) error {
	if len(lines) == 0 {
		return fmt.Errorf("no input lines")
	}
	header := lines[0]
	if !*caseSens {
		header = strings.ToLower(header)
	}
	existingFields := strings.Split(header, fieldSep)

	// save printable indices
	var colsToPrint = []field{}
	var nColsToPrint = 0
	for i, s := range existingFields {
		for _, rf := range requestedFields {
			if rf.normName == s {
				col := field{
					req:     rf,
					inIndex: i,
				}
				if *preserveOrder {
					col.outIndex = nColsToPrint
				} else {
					col.outIndex = rf.index
				}
				colsToPrint = append(colsToPrint, col)
				nColsToPrint++
			}
		}
	}

	// check for invalid columns in input flag
	for _, rf := range requestedFields {
		if !slices.Contains(existingFields, rf.normName) {
			return fmt.Errorf("requested field %s does not exist in input data", rf.name)
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
			for _, ff := range colsToPrint {
				if ff.inIndex == i {
					outFS[ff.outIndex] = f
				}
			}
		}
		outS := strings.Join(outFS, "\t")
		if li == 0 {
			if printedHeaders[outS] {
				continue
			}
			printedHeaders[outS] = true
			if len(printedHeaders) > 1 {
				return fmt.Errorf("Mismatching output headers: %s", strings.Join(maps.Keys(printedHeaders), "\n"))
			}
		}
		fmt.Println(outS)
	}
	return nil
}

const cmdname = "print_columns"

// options
var caseSens, skipHeader, preserveOrder *bool
var fieldSep string

var columnSplitRE = regexp.MustCompile("[,;: ]+")

func main() {

	caseSens = flag.Bool("c", false, "Case sensitive column headers")
	fieldSepFlag := flag.String("s", "<tab>", "Field `separator`")
	skipHeader = flag.Bool("H", false, "Skip output header")
	preserveFN := "o"
	preserveOrder = flag.Bool(preserveFN, false, "Preserve input file's column ordering")
	repeatFN := "r"
	allowRepeatedColumns := flag.Bool(repeatFN, false, "Allow repeated output fields")
	verb := flag.Bool("v", false, "Verbose output")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Script to selected columns based on file header")
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

	if flag.NArg() < 2 {
		printUsage()
		os.Exit(0)
	}

	fieldSep = *fieldSepFlag
	if fieldSep == "<tab>" {
		fieldSep = "\t"
	}

	if *preserveOrder && *allowRepeatedColumns {
		fmt.Fprintf(os.Stderr, "[warning] Flags -%s and -%s makes little sense to use in combination\n", preserveFN, repeatFN)
	}

	var requestedFieldsString = flag.Args()[0]
	var requestedFields = []reqField{}
	var seenRequestedFields = map[string]bool{}
	for i, name := range columnSplitRE.Split(requestedFieldsString, -1) {
		var key = name
		if !*caseSens {
			key = strings.ToLower(key)
		}
		f := reqField{index: i, normName: key, name: name}
		if _, repeated := seenRequestedFields[f.normName]; repeated && !*allowRepeatedColumns {
			fmt.Fprintf(os.Stderr, "[error] Repeated columns in request: %s (use flag -%s to allow repeated columns)\n\n", key, repeatFN)
			printUsage()
			os.Exit(1)
		}
		requestedFields = append(requestedFields, f)
		seenRequestedFields[f.normName] = true
	}

	if *verb {
		fmt.Fprintf(os.Stderr, "Separator: %#v\n", *fieldSepFlag)
		fmt.Fprintf(os.Stderr, "Requested fields: %s\n", requestedFieldsString)
	}

	if flag.NArg() > 1 {
		for _, f := range flag.Args()[1:] {
			if *verb {
				fmt.Fprintf(os.Stderr, "Reading file: %s\n", f)
			}
			lines, err := hio.ReadFileToLines(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[error] Failed to read from file %s: %v\n", f, err)
				os.Exit(1)
			}
			err = process(requestedFields, lines)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[error] Failed to process file %s: %v\n", f, err)
				os.Exit(1)
			}
		}
	} else { // read from stdin
		lines, err := hio.ReadStdinToLines()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[error] Failed to read from stdin: %v\n", err)
			os.Exit(1)
		}
		err = process(requestedFields, lines)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[error] Failed to process stdin: %v\n", err)
			os.Exit(1)
		}
	}

}
