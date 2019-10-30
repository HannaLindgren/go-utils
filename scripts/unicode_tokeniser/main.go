package main

import (
	"flag"
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"unicode"

	"github.com/HannaLindgren/go-utils/scripts/lib"
)

type interval struct {
	from int
	to   int
}

func (i interval) matches(n int) bool {
	return n >= i.from && n <= i.to
}

var numeric = interval{48, 57}
var spacing = regexp.MustCompile(`\s`)
var punctuation = regexp.MustCompile(`\pP`)

func blockFor(r rune) string {
	n := int(r)
	if numeric.matches(n) {
		return "Numeric"
	}
	if spacing.MatchString(string(r)) {
		return "White space"
	}
	if punctuation.MatchString(string(r)) {
		return "Punctuation"
	}
	for s, t := range unicode.Scripts {
		if unicode.In(r, t) {
			return s
		}
	}
	return "<UNDEF>"
}

func nfcNorm(s string) string {
	norm, _, _ := transform.String(norm.NFC, s)
	return norm
}

func nfdNorm(s string) string {
	norm, _, _ := transform.String(norm.NFD, s)
	return norm
}

func normalize(s string) string {
	if *nfc {
		return nfcNorm(s)
	} else if *nfd {
		return nfdNorm(s)
	}
	return s
}

func printAccu(accu []rune, block string) {
	if *xml {
		fmt.Printf("<token type='%s'>%s</token>\n", block, string(accu))
	} else {
		fmt.Printf("%s\t%s\n", string(accu), block)
	}
}

func process(s string) {
	sx := s
	var lastBlock string
	var accu = []rune{}
	for _, r := range []rune(normalize(sx)) {
		block := blockFor(r)
		if len(accu) > 0 && block != lastBlock {
			printAccu(accu, lastBlock)
			accu = []rune{}
		}
		accu = append(accu, r)
		lastBlock = block
	}
	if len(accu) > 0 {
		printAccu(accu, lastBlock)
	}
}

var nfc *bool
var nfd *bool
var xml *bool

func main() {
	cmdname := filepath.Base(os.Args[0])
	nfc = flag.Bool("c", false, "NFC -- Canonical composition on all input (default false)")
	nfd = flag.Bool("d", false, "NFD -- Canonical decomposition on all input (default false)")
	xml = flag.Bool("x", false, "XML output (default: tab-separated)")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Utility script to tokenise strings based on their unicode block.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       %s <string>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       echo <string> | %s\n", cmdname)
		// fmt.Fprintln(os.Stderr, cmdname+" <flags> <file1> <file2>")
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	if *nfd && *nfc {
		fmt.Fprintf(os.Stderr, "nfc and nfd options cannot be combined\n")
		printUsage()
		os.Exit(0)
	}

	if len(flag.Args()) > 0 {
		for _, arg := range os.Args[1:] {
			if lib.IsFile(arg) {
				text, err := lib.ReadFileToString(arg)
				if err != nil {
					log.Fatalf("%v", err)
				}
				process(text)
			} else {
				process(arg)
			}
		}
	} else {
		text, err := lib.ReadStdinToString()
		if err != nil {
			log.Fatalf("%v", err)
		}
		process(text)
	}
}
