package main

import (
	"fmt"
	"github.com/HannaLindgren/go-scripts/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func process(s string) string {
	runes := []rune(s)
	head := ""
	if len(runes) > 0 {
		head = strings.ToUpper(string(runes[0]))
	}
	tail := ""
	if len(runes) > 0 {
		tail = strings.ToLower(string(runes[1:]))
	}
	return head + tail
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		os.Exit(1)
	}
	err := util.ConvertFilesAndPrint(process, os.Args[1:])
	if err != nil {
		log.Fatalf("%v", err)
	}
}
