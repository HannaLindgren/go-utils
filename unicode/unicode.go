package main

import (
	"fmt"
	"github.com/HannaLindgren/go-scripts/util"
	"golang.org/x/text/unicode/runenames"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func getCharBlock(r rune) string {
	for s, t := range unicode.Scripts {
		if unicode.In(r, t) {
			return s
		}
	}
	return "<UNDEF>"
}
func process(s string) string {
	res := []string{}
	for _, r := range []rune(s) {
		name := runenames.Name(r)
		uc := fmt.Sprintf("%U", r)
		block := getCharBlock(r)
		res = append(res, fmt.Sprintf("%s\t%s\t%s\t%s", string(r), uc, name, block))
	}
	return strings.Join(res, "\n")
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-h") {
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
