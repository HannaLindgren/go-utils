package main

import (
	"fmt"
	"github.com/HannaLindgren/go-utils/scripts/util"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

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

func process(s string) string {
	res := []string{}
	for _, r := range []rune(s) {
		uc := codeFor(r)
		fmt.Printf("%s", uc)
	}
	return strings.Join(res, "\n")
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-h") {
		fmt.Fprintln(os.Stderr, "Utility script to convert strings into their unicode representation.")
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
