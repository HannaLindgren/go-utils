package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var ignoreCase *bool
var trim *bool
var verb *bool
var quiet *bool

const cmdname = "compare_line_by_line"

func readLines(file string) []string {
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Couldn't read file %s : %v", file, err)
	}
	lines := strings.Split(string(bts), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" { // sometimes I get an empty line last....
		lines = lines[:len(lines)-1]
	}
	return lines
}

func readLine(lines []string, lineNo int) (string, error) {
	if lineNo < len(lines) {
		return lines[lineNo], nil
	}
	return "", fmt.Errorf("Line %d is after EOF", lineNo)
}

func equal(s1, s2 string) bool {
	if *ignoreCase {
		return strings.EqualFold(s1, s2)
	}
	return s1 == s2
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func main() {
	ignoreCase = flag.Bool("i", false, "ignore case (default false)")
	trim = flag.Bool("t", false, "trim lines (default false)")
	verb = flag.Bool("v", false, "verbose -- print details for all lines, including those matching (default false)")
	quiet = flag.Bool("q", false, "quiet -- stats only, no line details (default false)")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, cmdname+" <flags> <file1> <file2>")
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	if *verb && *quiet {
		fmt.Fprintln(os.Stderr, "verbose and quiet flags cannot be combined!\n\nUsage:")
		printUsage()
		os.Exit(1)
	}

	if flag.NArg() != 2 {
		printUsage()
		os.Exit(0)
	}

	file1, file2 := flag.Arg(0), flag.Arg(1)
	if *verb {
		fmt.Fprintf(os.Stderr, "File1: %s\n", file1)
		fmt.Fprintf(os.Stderr, "File2: %s\n", file2)
	}
	if file1 == file2 {
		fmt.Printf("[%s] Comparing a file to itself doesn't make sense: %s\n", cmdname, file1)
		return
	}
	lines1, lines2 := readLines(file1), readLines(file2)
	n1, n2 := len(lines1), len(lines2)
	max := max(n1, n2)

	fmt.Fprintf(os.Stderr, "CaseSens:  %v\n", *ignoreCase)
	fmt.Fprintf(os.Stderr, "TrimSpace: %v\n", *trim)

	nMismatch, sizeDiff := 0, 0

	for i := 0; i < max; i++ {
		l1, eof1 := readLine(lines1, i)
		l2, eof2 := readLine(lines2, i)
		if eof1 != nil && eof2 == nil {
			sizeDiff++
			if !*quiet {
				fmt.Printf("F2 after F1\tL%d\t%s\n", i, l2)
			}
		} else if eof1 == nil && eof2 != nil {
			sizeDiff++
			if !*quiet {
				fmt.Printf("F1 after F2\tL%d\t%s\n", i, l1)
			}
		} else if equal(l1, l2) {
			if *verb {
				fmt.Printf("MATCH\tL%d\t%s\t%s\n", i, l1, l2)
			}
		} else {
			nMismatch++
			if !*quiet {
				fmt.Printf("DIFF\tL%d\t%s\t%s\n", i, l1, l2)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "F1: %s\n", file1)
	fmt.Fprintf(os.Stderr, "F2: %s\n", file2)

	fmt.Fprintf(os.Stderr, "F1 LINES READ:  %8d lines\n", n1)
	fmt.Fprintf(os.Stderr, "F2 LINES READ:  %8d lines\n", n2)
	fmt.Fprintf(os.Stderr, "FILE SIZE DIFF: %8d lines\n", sizeDiff)
	fmt.Fprintf(os.Stderr, "LINE DIFFS:     %8d lines\n", nMismatch)
}
