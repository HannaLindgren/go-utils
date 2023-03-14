package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/HannaLindgren/go-utils/scripts/lib"
	str "github.com/HannaLindgren/go-utils/strings"
	"github.com/HannaLindgren/go-utils/unicode"
)

var toker = unicode.Tokenizer{}

func convert(s string) string {
	res := ""
	for _, t := range toker.Tokenize(s) {
		res = res + str.UpcaseInitial(t.String, *downcaseRemainder)
	}
	if !strings.EqualFold(s, res) {
		panic(fmt.Sprintf("Expected output string to equal input string except for case, but found: <%s> => <%s>", s, res))
	}
	return res
}

var downcaseRemainder *bool

func main() {
	cmdname := "capitalize"
	downcaseRemainder = flag.Bool("d", false, "downcase remainder")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <flags> <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s <flags> \n", cmdname)
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	err := lib.ConvertAndPrintFromArgsOrStdin(convert, flag.Args())
	if err != nil {
		log.Fatalf("%v", err)
	}
	if err != nil {
		log.Fatalf("%v", err)
	}
}
