package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/HannaLindgren/go-utils/io"
)

var fieldSep = "\t"

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		os.Exit(1)
	}
	f := os.Args[1]
	r, fh, err := io.GetFileReader(f)
	defer fh.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}
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
