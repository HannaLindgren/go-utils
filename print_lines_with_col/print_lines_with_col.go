package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getFileReader(fName string) io.Reader {
	fs, err := os.Open(fName)
	if err != nil {
		log.Fatalf("Couldn't open file %s for reading : %v\n", fName, err)
	}

	if strings.HasSuffix(fName, ".gz") {
		gz, err := gzip.NewReader(fs)
		if err != nil {
			log.Fatalf("Couldn't to open gz reader : %v", err)
		}
		return io.Reader(gz)
	}
	return io.Reader(fs)
}

// TODO: Add flags for these settings
var ignoreCase *bool // = false
var fieldSep *string // = "\t"
var trimSpace *bool  // = false

var dbg = false

// static + dynamic variables
var lines = make(map[int]map[string][]string)
var indices = []int{}
var nPrinted = 0
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

func loadContentFile(fname string) {
	r := getFileReader(fname)
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		l := scan.Text()
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

func readFieldFile(fname string) {
	r := getFileReader(fname)
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		field0 := scan.Text()
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
					fmt.Println(line)
				}
				nPrinted++
			}
		}
		if !found {
			missing = append(missing, field0)
		}
	}
}

func debug(s string) {
	if dbg {
		fmt.Println(s)
	}
}

func debugf(template string, values ...interface{}) {
	debug(fmt.Sprintf(template, values))
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	ignoreCase = flag.Bool("i", false, "ignore case (default false)")
	trimSpace = flag.Bool("t", false, "trim lines (default false)")
	fieldSep = flag.String("f", "\t", "field separator")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Utility script to filter lines in a file based on a certain column, using a list of which field values to print.\n")
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, cmdname+" <input file> <field indices to check> <file with list of field values to print>")
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	if flag.NArg() != 3 {
		printUsage()
		os.Exit(0)
	}

	inputFile := flag.Arg(0)
	fields := flag.Arg(1)
	fieldsToPrint := flag.Arg(2)

	loadFieldIndices(fields)
	debugf("indices: %v\n", indices)

	loadContentFile(inputFile)
	debugf("%v\n", lines)

	readFieldFile(fieldsToPrint)

	fmt.Fprintf(os.Stderr, "PRINTED %d\n", nPrinted)
	fmt.Fprintf(os.Stderr, "MISSING %d\n", len(missing))

	if len(missing) > 0 {
		missingFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s.txt", cmdname))
		f, err := os.Create(missingFile)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
		w := bufio.NewWriter(f)
		for _, s := range missing {
			out := fmt.Sprintf("%s\n", s)
			_, err := w.WriteString(out)
			if err != nil {
				log.Fatal(err)
			}
		}
		w.Flush()

		fmt.Fprintf(os.Stderr, "MISSING PRINTED TO FILE %s\n", missingFile)
	}

}
