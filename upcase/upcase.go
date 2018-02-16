package main

import (
	"github.com/HannaLindgren/go-scripts/util"
	"os"
	"path/filepath"
	"strings"
)

func process(s string) string {
	return strings.ToUpper(s)
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	util.ConvertLinesFromStinOrFiles(cmdname, os.Args, process)
}
