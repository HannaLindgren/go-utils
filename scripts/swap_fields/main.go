package main

import (
	//"bufio"
	"fmt"
	"github.com/HannaLindgren/go-utils/scripts/lib"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var fieldSep = "\t"
var is = []int64{}

func process(s string) string {
	if s == "" {
		return s
	}
	fs := strings.Split(s, fieldSep)
	output := []string{}
	for _, i := range is {
		output = append(output, fs[i])
	}
	return strings.Join(output, fieldSep)
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		fmt.Println(len(os.Args))
		fmt.Fprintf(os.Stderr, "Usage: %s <fields-to-print> <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s <fields-to-print>\n", cmdname)
		os.Exit(1)
	}
	for _, f := range strings.Split(os.Args[1], ",") {
		i, err := strconv.ParseInt(strings.TrimSpace(f), 10, 64)
		if err != nil {
			log.Fatalf("Couldn't parse string to int : %v", err)
		}
		is = append(is, i-1)
	}

	args := []string{}
	if len(os.Args) >= 2 {
		args = os.Args[2:]
	}
	err := lib.ConvertAndPrintFromFilesOrStdin(process, args)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// for i := 2; i < len(os.Args); i++ {
	// 	f := os.Args[i]
	// 	r, fh, err := lib.GetFileReader(f)
	// 	defer fh.Close()
	// 	if err != nil {
	// 		log.Fatalf("%v", err)
	// 	}
	// 	scan := bufio.NewScanner(r)
	// 	for scan.Scan() {
	// 		s := scan.Text()
	// 		fs := strings.Split(s, fieldSep)
	// 		output := []string{}
	// 		for _, i := range is {
	// 			output = append(output, fs[i])
	// 		}
	// 		fmt.Println(strings.Join(output, fieldSep))
	// 	}
	// }
}
