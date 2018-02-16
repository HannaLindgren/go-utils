package main

import (
	"github.com/HannaLindgren/go-scripts/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func capitalize(s string) string {
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

var delims = strings.Split(" 	(){}\"”<>'-/&;.!?", "")

func tokenize(s string) []string {
	var acc = []string{s}
	for _, delim := range delims {
		tmp := []string{}
		for _, s := range acc {
			for _, split := range strings.SplitAfter(s, delim) {
				tmp = append(tmp, split)
			}
		}
		acc = tmp
	}
	return acc
}

func process(s string) string {
	res := ""
	for _, token := range tokenize(s) {
		res = res + capitalize(token)
	}
	if !strings.EqualFold(s, res) {
		log.Fatalf("Lost in conversion!\t%s\t%s", s, res)
	}
	return res
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	util.ConvertLinesFromStinOrFiles(cmdname, os.Args, process)
}
