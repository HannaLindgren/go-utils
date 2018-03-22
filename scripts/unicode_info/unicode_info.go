package main

import (
	"flag"
	"fmt"
	"github.com/HannaLindgren/go-utils/scripts/util"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/runenames"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ucNumberRe = regexp.MustCompile(`^(?:\\u|[uU][+])([a-fA-F0-9]{4})$`)

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

func normalize(s string) string {
	if *nfc {
		return string(norm.NFC.Bytes([]byte(s)))
	} else if *nfd {
		return string(norm.NFD.Bytes([]byte(s)))
	}
	return s
}

func process(s string) string {
	res := []string{}
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
	for _, r := range []rune(normalize(sx)) {
		name := runenames.Name(r)
		uc := codeFor(r)
		block := blockFor(r)
		res = append(res, fmt.Sprintf("%s\t%s\t%s\t%s", string(r), uc, name, block))
	}
	return strings.Join(res, "\n")
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

	err := util.ConvertAndPrintFromFilesOrStdin(process, flag.Args())
	if err != nil {
		log.Fatalf("%v", err)
	}
}
