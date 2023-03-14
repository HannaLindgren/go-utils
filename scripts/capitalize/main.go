package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/HannaLindgren/go-utils/scripts/lib"
	str "github.com/HannaLindgren/go-utils/strings"
	"github.com/HannaLindgren/go-utils/unicode"
)

var toker = unicode.Tokenizer{}

func convert(s string) string {
	res := ""
	for _, t := range toker.Tokenize(s) {
		res = res + str.UpcaseInitial(t.String, false)
	}
	if !strings.EqualFold(s, res) {
		panic(fmt.Sprintf("Expected output string to equal input string except for case, but found: <%s> => <%s>", s, res))
	}
	return res
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-h") {
		fmt.Fprintf(os.Stderr, "Capitalize each word (as separated by a very simple tokenizer)")
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s\n", cmdname)
		os.Exit(1)
	}
	err := lib.ConvertAndPrintFromArgsOrStdin(convert, os.Args[1:])
	if err != nil {
		log.Fatalf("%v", err)
	}
}
