package main

import (
	"fmt"
	"github.com/HannaLindgren/go-utils/scripts/lib"
	"golang.org/x/text/unicode/norm"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func process(s string) string {
	return string(norm.NFC.Bytes([]byte(s)))
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-h") {
		fmt.Fprintln(os.Stderr, "Utility script for canonical composition of text data (non-destructive conversion).")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s\n", cmdname)
		os.Exit(1)
	}
	err := lib.ConvertAndPrintFromFileArgsOrStdin(process)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
