package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/HannaLindgren/go-utils/io"
	"github.com/HannaLindgren/go-utils/unicode"
)

var up unicode.Processor

func process(s string) {
	for _, ui := range up.UnicodeInfo(s) {
		fmt.Printf("%s\t%s\t%s\t%s\n", ui.String, ui.Unicode, ui.CharName, ui.CodeBlock)
	}
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	convertFromUnicodeNumbers := flag.Bool("u", false, "unicode -- convert from unicode numbers (default false)")
	nfc := flag.Bool("c", false, "NFC -- Canonical composition on all input (default false)")
	nfd := flag.Bool("d", false, "NFD -- Canonical decomposition on all input (default false)")

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

	up = unicode.Processor{
		NFC:                       *nfc,
		NFD:                       *nfd,
		ConvertFromUnicodeNumbers: *convertFromUnicodeNumbers,
	}

	if len(flag.Args()) > 0 {
		for _, arg := range os.Args[1:] {
			if io.IsFile(arg) {
				text, err := io.ReadFileToString(arg)
				if err != nil {
					log.Fatalf("%v", err)
				}
				process(text)
			} else {
				process(arg)
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
