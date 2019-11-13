package main

import (
	"flag"
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/runenames"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/HannaLindgren/go-utils/scripts/lib"
)

var ucNumberRe = regexp.MustCompile(`^(?:\\u|[uU][+])([a-fA-F0-9]{4})$`)

const newline rune = '\n'

var hardwiredNames = map[rune]string{
	newline: "NEWLINE",
	'	': "TAB",
}

func inhibitSpecialChar(r rune) bool {
	_, ok := hardwiredNames[r]
	return ok
}

func isUcNumber(s string) bool {
	return ucNumberRe.MatchString(s)
}

func ucNumber2String(s string) string {
	is := ucNumberRe.ReplaceAllString(s, "$1")
	i, err := strconv.ParseInt(is, 16, 32)
	if err != nil {
		log.Fatalf("Couldn't parse unicode number from input string '%s' : %v", s, err)
	}
	r := rune(i)
	return string(r)
}

func blockFor(r rune) string {
	for s, t := range unicode.Scripts {
		if unicode.In(r, t) {
			return s
		}
	}
	return "<UNDEF>"
}

func codeFor(r rune) string {
	uc := fmt.Sprintf("%U", r)
	return fmt.Sprintf("\\u%s", uc[2:])
}

func nameFor(r rune) string {
	if name, ok := hardwiredNames[r]; ok {
		return name
	}
	return runenames.Name(r)
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
	if *convertFromUnicodeNumbers {
		tmp := []string{}
		for _, chunk := range strings.Split(sx, " ") {
			if isUcNumber(chunk) {
				tmp = append(tmp, ucNumber2String(chunk))
			} else {
				tmp = append(tmp, chunk)
			}
		}
		sx = strings.Join(tmp, " ")
	}
	for _, r := range normalize(sx) {
		name := nameFor(r)
		uc := codeFor(r)
		block := blockFor(r)
		thisS := string(r)
		if inhibitSpecialChar(r) {
			thisS = ""
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", thisS, uc, name, block)
	}
}

var nfc *bool
var nfd *bool
var convertFromUnicodeNumbers *bool

func main() {
	cmdname := filepath.Base(os.Args[0])
	convertFromUnicodeNumbers = flag.Bool("u", false, "unicode -- convert from unicode numbers (default false)")
	nfc = flag.Bool("c", false, "NFC -- Canonical composition on all input (default false)")
	nfd = flag.Bool("d", false, "NFD -- Canonical decomposition on all input (default false)")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Utility script to retrieve information about input characters: unicode number, unicode name, and character block.")
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
