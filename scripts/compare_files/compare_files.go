package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var f1Notf2 = 0
var f2Notf1 = 0
var nDiff = 0
var nBoth = 0
var nLines1 = 0
var nLines2 = 0
var sizeDiff = 0

type output int

const (
	f1 output = iota
	f2
	all
	both
	diff
	stats
)

var modes = []output{f1, f2, all, both, diff, stats}
var modesString string = strings.Join(strings.Fields(fmt.Sprint(modes)), "|")

func modesHelp(prefix string) string {
	return strings.Join([]string{
		prefix + f1.String() + ":    Lines in file1 only",
		prefix + f2.String() + ":    Lines in file2 only",
		prefix + all.String() + ":   All lines (with diff info)",
		prefix + both.String() + ":  Lines occurring in both files",
		prefix + diff.String() + ":  Mismatching lines",
		prefix + stats.String() + ": Statistics",
	}, "\n")
}

const _outputName = "f1f2allbothdiffstats"

var _outputIndex = [...]uint8{0, 2, 4, 7, 11, 15, 20}

func string2output(s string) output {
	switch s {
	case "f1":
		return f1
	case "f2":
		return f2
	case "all":
		return all
	case "both":
		return both
	case "diff":
		return diff
	case "stats":
		return stats
	}
	log.Fatalf("Invalid output mode: %s", s)
	return stats
}
func (i output) String() string {
	if i < 0 || i >= output(len(_outputIndex)-1) {
		return fmt.Sprintf("output(%d)", i)
	}
	return _outputName[_outputIndex[i]:_outputIndex[i+1]]
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
	var finalNewline = regexp.MustCompile("\n$")
	var s = finalNewline.ReplaceAllString(string(bts), "")
	return strings.Split(s, "\n")
}

func readLine(lines []string, lineNo int) (string, error) {
	if lineNo < len(lines) {
		return lines[lineNo], nil
	}
	return "", fmt.Errorf("Line %d is after EOF", lineNo)
}

var ignoreCase *bool
var keepOrdering *bool
var defaultMode = stats
var mode output
var trim *bool

func equal(s1, s2 string) bool {
	if *ignoreCase {
		return strings.EqualFold(s1, s2)
	}
	return s1 == s2
}

func unsorted(file1, file2 string) {
	lines1, lines2 := readLines(file1), readLines(file2)
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
			if mode == both {
				fmt.Println(l0)
			} else if mode == all {
				fmt.Printf("f1 & f2\t%s\n", l0)
			}
		} else {
			f2Notf1++
			nDiff++
			if mode == f2 {
				fmt.Println(l)
			} else if mode == all || mode == diff {
				fmt.Printf("f2 not f1\t%s\n", l)
			}
		}
	}
	for _, inputs := range lines {
		for _, input := range inputs {
			if _, ok := found[input]; !ok {
				f1Notf2++
				if mode == f1 {
					fmt.Println(input)
				} else if mode == all || mode == diff {
					fmt.Printf("f1 not f2\t%s\n", input)
				}
			}
		}
	}
}

func lineByLine(file1, file2 string) {
	lines1, lines2 := readLines(file1), readLines(file2)
	nLines1, nLines2 = len(lines1), len(lines2)
	max := max(nLines1, nLines2)

	for i := 0; i < max; i++ {
		l1, eof1 := readLine(lines1, i)
		l2, eof2 := readLine(lines2, i)
		if eof1 != nil && eof2 == nil {
			nDiff++
			sizeDiff++
			if mode == all || mode == diff || mode == f1 {
				fmt.Printf("f2 after f1\tL%d\t%s\n", i, l2)
			}
		} else if eof1 == nil && eof2 != nil {
			nDiff++
			sizeDiff++
			if mode == all || mode == diff || mode == f2 {
				fmt.Printf("f1 after f2\tL%d\t%s\n", i, l1)
			}
		} else if equal(l1, l2) {
			nBoth++
			if mode == both {
				fmt.Println(l1)
			} else if mode == all {
				fmt.Printf("f1 & f2\tL%d\t%s\n", i, l2)
			}
		} else {
			f1Notf2++
			f2Notf1++
			nDiff++
			if mode == f1 {
				fmt.Println(l1)
			}
			if mode == f2 {
				fmt.Println(l2)
			}
			if mode == all || mode == diff {
				fmt.Printf("f1 not f2\tL%d\t%s\n", i, l1)
				fmt.Printf("f2 not f1\tL%d\t%s\n", i, l2)
			}
		}
	}

}

func internalInitTests() {
	for _, o := range modes {
		s := o.String()
		o2 := string2output(s)
		if o2 != o {
			log.Fatalf("Internal init error for output type: %s <=> %s", o, o2)
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
		mode = string2output(*modeF)
	}

	fmt.Fprintf(os.Stderr, "IgnoreCase:   %v\n", *ignoreCase)
	fmt.Fprintf(os.Stderr, "KeepOrdering: %v\n", *keepOrdering)
	fmt.Fprintf(os.Stderr, "TrimSpace:    %v\n", *trim)
	fmt.Fprintf(os.Stderr, "Mode:         %s\n", mode.String())

	if *keepOrdering {
		lineByLine(file1, file2)
	} else {
		unsorted(file1, file2)
	}

	if mode == stats {
		fmt.Printf("F1: %s\n", file1)
		fmt.Printf("F2: %s\n", file2)
		fmt.Printf("F1 LINES READ:  %8d lines\n", nLines1)
		fmt.Printf("F2 LINES READ:  %8d lines\n", nLines2)
		fmt.Printf("FILE SIZE DIFF: %8d lines\n", sizeDiff)
		fmt.Printf("F1 not F2       %8d lines\n", f1Notf2)
		fmt.Printf("F2 not F1       %8d lines\n", f2Notf1)
		fmt.Printf("F1  &  F2       %8d lines\n", nBoth)
		fmt.Printf("TOTAL DIFF      %8d lines\n", nDiff)
	}
}
