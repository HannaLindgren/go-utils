package main

import (
	"fmt"
	"github.com/HannaLindgren/go-utils/scripts/lib"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var char string

func process(s string) string {
	var res = false
	if strings.Contains(s, char) {
		res = true
	}
	return fmt.Sprintf("%s\t%v", s, res)
}

func code2char(s string) (string, error) {
	i, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		return "", err
	}
	r := rune(i)
	return string(r), nil
}

func main() {
	var err error
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-h") {
		fmt.Fprintln(os.Stderr, "Utility script to search for unicode characters.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintf(os.Stderr, "Usage: %s <char|charcode> <files>\n", cmdname)
		os.Exit(1)
	}
	char = os.Args[1]
	if strings.HasPrefix(strings.ToLower(char), "\\u") {
		char = strings.Replace(strings.ToLower(char), "\\u", "", -1)
		char, err = code2char(char)
		if err != nil {
			log.Fatalf("%v", err)
		}
	} else if strings.HasPrefix(strings.ToLower(char), "u") {
		char = strings.Replace(strings.ToLower(char), "u", "", -1)
		char, err = code2char(char)
		if err != nil {
			log.Fatalf("%v", err)
		}
	} else if strings.HasPrefix(strings.ToLower(char), "u+") {
		char = strings.Replace(strings.ToLower(char), "u", "", -1)
		char, err = code2char(char)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	err = lib.ConvertAndPrintFromFilesOrStdin(process, os.Args[2:])
	if err != nil {
		log.Fatalf("%v", err)
	}
}
