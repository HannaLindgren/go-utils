package main

import (
	"github.com/HannaLindgren/go-scripts/util"
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
	util.ConvertLinesFromStinOrFiles(cmdname, os.Args, process)
}
