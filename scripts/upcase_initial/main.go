package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/HannaLindgren/go-utils/scripts/lib"
	"github.com/HannaLindgren/go-utils/strings"
)

func convert(s string) string {
	return strings.UpcaseInitial(s, *downcaseRemainder)
}

var downcaseRemainder *bool

func main() {
	cmdname := "upcase_initial"
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
}
