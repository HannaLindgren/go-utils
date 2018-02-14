package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <fields-to-print> <files>\n", cmdname)
		os.Exit(1)
	}
	is := []int64{}
	for _, f := range strings.Split(os.Args[1], ",") {
		i, err := strconv.ParseInt(strings.TrimSpace(f), 10, 64)
		if err != nil {
			log.Fatalf("Couldn't parse string to int : %v", err)
		}
		is = append(is, i-1)
	}

	for i := 2; i < len(os.Args); i++ {
		f := os.Args[i]
		r := getFileReader(f)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			s := scan.Text()
			fs := strings.Split(s, fieldSep)
			output := []string{}
			for _, i := range is {
				output = append(output, fs[i])
			}
			fmt.Println(strings.Join(output, fieldSep))
		}
	}
}
