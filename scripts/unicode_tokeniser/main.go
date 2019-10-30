package main

import (
	"flag"
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"os"
	"path/filepath"
	"unicode"

	"github.com/HannaLindgren/go-utils/scripts/lib"
)

var numerals = map[rune]bool{
	'1': true, '2': true, '3': true, '4': true, '5': true, '6': true, '7': true, '8': true, '9': true, '0': true,
}

func blockFor(r rune) string {
	if _, ok := numerals[r]; ok {
		return "Numeric"
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

func process(s string) {
	sx := s
	var lastBlock string
	for _, r := range []rune(normalize(sx)) {
		block := blockFor(r)
		if lastBlock != "" && block != lastBlock {
			fmt.Printf("\t%s\n", lastBlock)
		}
		fmt.Print(string(r))
		lastBlock = block
	}
	fmt.Printf("\t%s\n", lastBlock)
}

var nfc *bool
var nfd *bool

func main() {
	cmdname := filepath.Base(os.Args[0])
	nfc = flag.Bool("c", false, "NFC -- Canonical composition on all input (default false)")
	nfd = flag.Bool("d", false, "NFD -- Canonical decomposition on all input (default false)")

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
