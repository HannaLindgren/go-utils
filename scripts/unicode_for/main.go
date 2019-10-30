package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/HannaLindgren/go-utils/scripts/lib"
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

func processRune(r rune) string {
	uc := codeFor(r)
	return fmt.Sprintf("%s", uc)
}

const newline rune = '\n'

func process(s string) {
	for _, r := range []rune(s) {
		fmt.Print(processRune(r))
		if r == newline {
			fmt.Println()
		}
	}
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
	if len(os.Args[1:]) > 0 {
		for _, arg := range os.Args[1:] {
			if lib.IsFile(arg) {
				text, err := lib.ReadFileToString(arg)
				if err != nil {
					log.Fatalf("%v", err)
				}
				process(text)
			} else {
				process(arg)
				fmt.Println()
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
