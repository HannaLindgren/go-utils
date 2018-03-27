package main

import (
	"fmt"
	"github.com/HannaLindgren/go-utils/scripts/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func process(s string) string {
	runes := []rune(s)
	len := len(runes)
	if len > 0 {
		return fmt.Sprintf("%d\t%s", len, s)
	}
	return ""
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-h") {
		fmt.Fprintf(os.Stderr, "Prints the length of each input line\n")
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s\n", cmdname)
		os.Exit(1)
	}
	err := util.ConvertAndPrintFromFileArgsOrStdin(process)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
