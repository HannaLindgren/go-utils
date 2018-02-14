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

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		os.Exit(1)
	}
	for i := 1; i < len(os.Args); i++ {
		f := os.Args[i]
		r := getFileReader(f)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			s := scan.Text()
			res := capitalize(s)
			fmt.Println(res)
		}
	}
}
