package main

import (
	//"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/HannaLindgren/go-utils/scripts/lib"
	"github.com/HannaLindgren/go-utils/strings"
)

var fieldSep = "\t"
var is = []int64{}

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) < 1 {
		fmt.Fprintln(os.Stderr, "Reverse each line")
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s\n", cmdname)
		os.Exit(1)
	}
	err := lib.ConvertAndPrintFromFilesOrStdin(strings.Reverse, os.Args[1:])
	if err != nil {
		log.Fatalf("%v", err)
	}

}
