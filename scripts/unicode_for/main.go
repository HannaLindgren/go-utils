package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/HannaLindgren/go-utils/io"
	"github.com/HannaLindgren/go-utils/unicode"
)

func process(s string) {
	for _, r := range s {
		fmt.Print(unicode.UnicodeForR(r))
		if r == unicode.Newline {
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
			if io.IsFile(arg) {
				text, err := io.ReadFileToString(arg)
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
		text, err := io.ReadStdinToString()
		if err != nil {
			log.Fatalf("%v", err)
		}
		process(text)
	}
}
