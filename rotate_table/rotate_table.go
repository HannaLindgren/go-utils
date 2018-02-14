package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getFileReader(fName string) io.Reader {
	fs, err := os.Open(fName)
	if err != nil {
		log.Fatalf("Couldn't open file %s for reading : %v\n", fName, err)
	}

	if strings.HasSuffix(fName, ".gz") {
		gz, err := gzip.NewReader(fs)
		if err != nil {
			log.Fatalf("Couldn't to open gz reader : %v", err)
		}
		return io.Reader(gz)
	}
	return io.Reader(fs)
}

var fieldSep = "\t"

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		os.Exit(1)
	}
	f := os.Args[1]
	r := getFileReader(f)
	scan := bufio.NewScanner(r)
	rows := [][]string{}
	for scan.Scan() {
		s := scan.Text()
		fs := strings.Split(s, fieldSep)
		for i, field := range fs {
			for len(rows) <= i {
				rows = append(rows, []string{})
			}
			rows[i] = append(rows[i], field)
		}
	}
	for _, row := range rows {
		fmt.Println(strings.Join(row, fieldSep))
	}
}
