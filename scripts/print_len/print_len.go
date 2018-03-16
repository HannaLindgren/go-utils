package main

import (
	"bufio"
	"fmt"
	"github.com/HannaLindgren/go-utils/scripts/util"
	"log"
	"os"
	"path/filepath"
)

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Prints the length of each input line\nUsage: %s <files>\n", cmdname)
		os.Exit(1)
	}
	for i := 1; i < len(os.Args); i++ {
		f := os.Args[i]
		r, fh, err := util.GetFileReader(f)
		defer fh.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			s := scan.Text()
			runes := []rune(s)
			len := len(runes)
			if len > 0 {
				fmt.Printf("%d\t%s\n", len, s)
			}
		}
	}
}
