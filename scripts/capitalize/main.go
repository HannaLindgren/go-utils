package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/HannaLindgren/go-utils/scripts/lib"
	str "github.com/HannaLindgren/go-utils/strings"
)

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-h") {
		fmt.Fprintf(os.Stderr, "Capitalize each word (as separated by a very simple tokenizer)")
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s\n", cmdname)
		os.Exit(1)
	}
	err := lib.ConvertAndPrintFromArgsOrStdin(str.CapitalizeTokens, os.Args[1:])
	if err != nil {
		log.Fatalf("%v", err)
	}
}
