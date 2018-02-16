package main

import (
	"bufio"
	"fmt"
	"github.com/HannaLindgren/go-scripts/util"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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
		r, fh, err := util.GetFileReader(f)
		defer fh.Close()
		if err != nil {
			log.Fatalf("%v", err)
		}
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
