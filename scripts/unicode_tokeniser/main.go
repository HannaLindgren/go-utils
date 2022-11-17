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

func main() {
	cmdname := filepath.Base(os.Args[0])
	nfc := flag.Bool("c", false, "NFC -- Canonical composition on all input (default false)")
	nfd := flag.Bool("d", false, "NFD -- Canonical decomposition on all input (default false)")
	xml := flag.Bool("x", false, "XML output (default: tab-separated)")

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

	toker := unicode.Tokenizer{
		UP: unicode.Processor{
			NFC: *nfc,
			NFD: *nfd,
		},
	}

	var process = func(s string) {
		for _, t := range toker.Tokenize(s) {
			if *xml {
				fmt.Printf("<token type='%s'>%s</token>\n", t.UnicodeBlock, t.String)
			} else {
				fmt.Printf("%s\t%s\n", t.String, t.UnicodeBlock)
			}
		}
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
