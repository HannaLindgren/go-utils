package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HannaLindgren/go-utils/io"
)

var fieldSep *string // = "\t"
var ignoreCase,
	printMissing,
	trimSpace *bool

// static/dynamic variables
var lines = make(map[int]map[string][]string)
var indices = []int{}
var nPrinted = 0
var nFound = 0
var missing = []string{}

func loadFieldIndices(fields string) {
	for _, s := range strings.Split(fields, ",") {
		i0, err := strconv.ParseInt(s, 10, 64)
		i := int(i0 - 1)
		if err != nil {
			log.Fatalf("Couldn't parse field index <%s> in input definition <%s>", s, fields)
		}
		indices = append(indices, i)
		lines[i] = make(map[string][]string)
	}
}

func loadContents(inputFile *string) {
	var lns []string
	if inputFile != nil {
		r, fh, err := io.GetFileReader(*inputFile)
		defer fh.Close()
		if err != nil {
			log.Fatalf("Couldn't read from content file %s: %v", *inputFile, err)
		}
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			lns = append(lns, scan.Text())

		}
	} else {
		stdin := bufio.NewReader(os.Stdin)
		b, err := ioutil.ReadAll(stdin)
		if err != nil {
			log.Fatalf("Couldn't read contents from stdin: %v", err)
		}
		lns = strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
	}
	for _, l := range lns {
		if *trimSpace {
			l = strings.TrimSpace(l)
		}
		for i, f := range strings.Split(l, *fieldSep) {
			if *ignoreCase {
				f = strings.ToUpper(f)
			}
			if _, ok := lines[i]; !ok {
				lines[i] = make(map[string][]string)
				lines[i][f] = []string{}
			}
			lines[i][f] = append(lines[i][f], l)
		}
	}
}

func readFields(fNameOrString string) {
	var fields []string
	if _, err := os.Stat(fNameOrString); errors.Is(err, os.ErrNotExist) {
		fields = append(fields, fNameOrString)
	} else {
		r, fh, err := io.GetFileReader(fNameOrString)
		defer fh.Close()
		if err != nil {
			log.Fatalf("Couldn't read field file %s: %v", fNameOrString, err)
		}
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			fields = append(fields, scan.Text())
		}
	}
	for _, field0 := range fields {
		field := field0
		if *ignoreCase {
			field = strings.ToUpper(field)
		}
		if *trimSpace {
			field = strings.TrimSpace(field)
		}
		found := false
		for _, i := range indices {
			if val, ok := lines[i][field]; ok {
				for _, line := range val {
					found = true
					nPrinted++
					if !*printMissing {
						fmt.Println(line)
					}
				}
				nFound++
			}
		}
		if !found {
			missing = append(missing, field0)
		}
	}
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	ignoreCase = flag.Bool("i", false, "ignore case (default false)")
	trimSpace = flag.Bool("t", false, "trim lines (default false)")
	printMissing = flag.Bool("m", false, "print missing items only (default false)")
	fieldSep = flag.String("f", "\t", "field separator")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Utility script to filter lines in a file based on a certain column, using a list of which field values to print.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, cmdname+" <input file> <field indices to check> <file with list of field values to print>")
		fmt.Fprintln(os.Stderr, " OR")
		fmt.Fprintln(os.Stderr, "cat <input> | "+cmdname+" <field indices to check> <file with list of field values to print>")
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	var inputFile *string
	var fields, fieldsToPrint string
	if flag.NArg() == 2 {
		fields = flag.Arg(0)
		fieldsToPrint = flag.Arg(1)
	} else if flag.NArg() == 3 {
		inpF := flag.Arg(0)
		inputFile = &inpF
		fields = flag.Arg(1)
		fieldsToPrint = flag.Arg(2)
	} else {
		printUsage()
		os.Exit(0)
	}

	loadFieldIndices(fields)
	loadContents(inputFile)
	readFields(fieldsToPrint)

	var foundNotPrinted string
	if *printMissing && nFound > 0 {
		foundNotPrinted = " *** NOT PRINTED"
	}
	fmt.Fprintf(os.Stderr, "Found %d entries/%d lines%s\n", nFound, nPrinted, foundNotPrinted)
	fmt.Fprintf(os.Stderr, "Missing entries: %d\n", len(missing))

	if *printMissing && len(missing) > 0 {
		for _, s := range missing {
			fmt.Printf("%s\n", s)
		}
	}

}
