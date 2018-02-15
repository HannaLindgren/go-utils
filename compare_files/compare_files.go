package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

var f1NotF2 = 0
var f2NotF1 = 0
var nDiff = 0
var nBoth = 0
var nLines1 = 0
var nLines2 = 0
var sizeDiff = 0

type Output int

const (
	F1 Output = iota
	F2
	All
	Both
	Diff
	Stats
)

var modes = []Output{F1, F2, All, Both, Diff, Stats}
var modesString string = strings.Join(strings.Fields(fmt.Sprint(modes)), "|")

func modesHelp(prefix string) string {
	return strings.Join([]string{
		prefix + F1.String() + ":    Lines in file1 only",
		prefix + F2.String() + ":    Lines in file2 only",
		prefix + All.String() + ":   All lines (with diff info)",
		prefix + Both.String() + ":  Lines occuring in both files",
		prefix + Diff.String() + ":  Mismatching lines",
		prefix + Stats.String() + ": Statistics",
	}, "\n")
}

const _Output_name = "F1F2AllBothDiffStats"

var _Output_index = [...]uint8{0, 2, 4, 7, 11, 15, 20}

func String2Output(s string) Output {
	switch s {
	case "F1":
		return F1
	case "F2":
		return F2
	case "All":
		return All
	case "Both":
		return Both
	case "Diff":
		return Diff
	case "Stats":
		return Stats
	}
	log.Fatalf("Invalid output mode: %s", s)
	return Stats
}
func (i Output) String() string {
	if i < 0 || i >= Output(len(_Output_index)-1) {
		return fmt.Sprintf("Output(%d)", i)
	}
	return _Output_name[_Output_index[i]:_Output_index[i+1]]
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func readLines(file string) []string {
	bts, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Couldn't read file %s : %v", file, err)
	}
	return strings.Split(string(bts), "\n")
}

func readLine(lines []string, lineNo int) (string, error) {
	if lineNo < len(lines) {
		return lines[lineNo], nil
	}
	return "", fmt.Errorf("Line %d is after EOF", lineNo)
}

var ignoreCase *bool
var keepOrdering *bool
var defaultMode Output = Stats
var mode Output
var trim *bool

func equal(s1, s2 string) bool {
	if *ignoreCase {
		return strings.EqualFold(s1, s2)
	}
	return s1 == s2
}

func unsorted(f1, f2 string) {
	lines1, lines2 := readLines(f1), readLines(f2)
	nLines1, nLines2 = len(lines1), len(lines2)
	sizeDiff = int(math.Abs(float64(nLines2 - nLines1)))
	lines := make(map[string][]string)
	found := make(map[string]bool)
	for _, l0 := range lines1 {
		l := l0
		if *ignoreCase {
			l = strings.ToLower(l0)
		}
		lines[l] = append(lines[l], l0)
	}
	for _, l0 := range lines2 {
		l := l0
		if *ignoreCase {
			l = strings.ToLower(l0)
		}
		inputs, exists := lines[l]
		if exists {
			nBoth++
			for _, input := range inputs {
				found[input] = true
			}
			if mode == Both {
				fmt.Println(l0)
			} else if mode == All {
				fmt.Printf("F1 & F2\t%s\n", l0)
			}
		} else {
			f2NotF1++
			nDiff++
			if mode == F2 {
				fmt.Println(l)
			} else if mode == All || mode == Diff {
				fmt.Printf("F2 not F1\t%s\n", l)
			}
		}
	}
	for _, inputs := range lines {
		for _, input := range inputs {
			if _, ok := found[input]; !ok {
				f1NotF2++
				if mode == F1 {
					fmt.Println(input)
				} else if mode == All || mode == Diff {
					fmt.Printf("F1 not F2\t%s\n", input)
				}
			}
		}
	}
}

func lineByLine(f1, f2 string) {
	lines1, lines2 := readLines(f1), readLines(f2)
	nLines1, nLines2 = len(lines1), len(lines2)
	max := max(nLines1, nLines2)

	for i := 0; i < max; i++ {
		l1, eof1 := readLine(lines1, i)
		l2, eof2 := readLine(lines2, i)
		if eof1 != nil && eof2 == nil {
			nDiff++
			sizeDiff++
			if mode == All || mode == Diff || mode == F1 {
				fmt.Printf("F1 after F2\tL%d\t%s\n", i, l1)
			}
		} else if eof1 == nil && eof2 != nil {
			nDiff++
			sizeDiff++
			if mode == All || mode == Diff || mode == F2 {
				fmt.Printf("F2 after F1\tL%d\t%s\n", i, l2)
			}
		} else if equal(l1, l2) {
			nBoth++
			if mode == Both {
				fmt.Println(l1)
			} else if mode == All {
				fmt.Printf("F1 & F2\tL%d\t%s\n", i, l2)
			}
		} else {
			f1NotF2++
			f2NotF1++
			nDiff++
			if mode == F1 {
				fmt.Println(l1)
			}
			if mode == F2 {
				fmt.Println(l2)
			}
			if mode == All || mode == Diff {
				fmt.Printf("F1 not F2\tL%d\t%s\n", i, l1)
				fmt.Printf("F2 not F1\tL%d\t%s\n", i, l2)
			}
		}
	}

}

func internalInitTests() {
	for _, o := range modes {
		s := o.String()
		o2 := String2Output(s)
		if o2 != o {
			log.Fatalf("Internal init error for Output type: %s <=> %s", o, o2)
		}
	}
}

func main() {

	cmdname := filepath.Base(os.Args[0])

	internalInitTests()

	ignoreCase = flag.Bool("i", false, "ignore case (default false)")
	keepOrdering = flag.Bool("o", false, "keep line ordering (default false)")
	trim = flag.Bool("t", false, "trim lines (default false)")
	var modeF = flag.String("m", "", fmt.Sprintf("output mode (default %s)\n%s\n         ", defaultMode, modesHelp("          ")))

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

	if flag.NArg() != 2 {
		printUsage()
		os.Exit(0)
	}

	file1, file2 := flag.Arg(0), flag.Arg(1)
	if file1 == file2 {
		fmt.Printf("[%s] Comparing a file to itself doesn't make sense: %s\n", cmdname, file1)
		return
	}

	if *modeF == "" {
		mode = defaultMode
	} else {
		mode = String2Output(*modeF)
	}

	fmt.Fprintf(os.Stderr, "CaseSens:     %v\n", *ignoreCase)
	fmt.Fprintf(os.Stderr, "KeepOrdering: %v\n", *keepOrdering)
	fmt.Fprintf(os.Stderr, "TrimSpace:    %v\n", *trim)
	fmt.Fprintf(os.Stderr, "Mode:         %s\n", mode.String())

	if *keepOrdering {
		lineByLine(file1, file2)
	} else {
		unsorted(file1, file2)
	}

	if mode == Stats {
		fmt.Printf("F1 LINES READ:  %8d lines\n", nLines1)
		fmt.Printf("F2 LINES READ:  %8d lines\n", nLines2)
		fmt.Printf("FILE SIZE DIFF: %8d lines\n", sizeDiff)
		fmt.Printf("F1 not F2       %8d lines\n", f1NotF2)
		fmt.Printf("F2 not F1       %8d lines\n", f2NotF1)
		fmt.Printf("F1  &  F2       %8d lines\n", nBoth)
		fmt.Printf("TOTAL DIFF      %8d lines\n", nDiff)
	}
}
