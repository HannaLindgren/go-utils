package main

import (
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

var ucNumberRe = regexp.MustCompile(`^(\\u|U[+])[a-f0-9]{4}$`)

func isUcNumber(s string) bool {
	return ucNumberRe.MatchString(s)
}

func ucNumber2String(s string) string {
	s = strings.Replace(s, `\u`, "", -1)
	s = strings.Replace(s, `U+`, "", -1)
	i, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		log.Fatal(err)
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
	if nfc {
		return string(norm.NFC.Bytes([]byte(s)))
	}
	return s
}

func process(s string) string {
	res := []string{}
	sx := s
	if isUcNumber(s) {
		sx = ucNumber2String(s)
	}
	for _, r := range []rune(normalize(sx)) {
		name := runenames.Name(r)
		uc := codeFor(r)
		block := blockFor(r)
		res = append(res, fmt.Sprintf("%s\t%s\t%s\t%s", string(r), uc, name, block))
	}
	return strings.Join(res, "\n")
}

var nfc = false

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-h") {
		fmt.Fprintln(os.Stderr, "Utility script to retrive information about input characters: unicode number, unicode name, and character block.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s\n", cmdname)
		os.Exit(1)
	}
	err := util.ConvertAndPrintFromFileArgsOrStdin(process)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
