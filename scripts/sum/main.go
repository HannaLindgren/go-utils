package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	hio "github.com/HannaLindgren/go-utils/io"
)

func main() {
	sum := 0.0
	n := 0
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-h") {
		fmt.Fprintf(os.Stderr, "Usage: %s <files> (one float/integer per line)\n", cmdname)
		fmt.Fprintln(os.Stderr, "       OR")
		fmt.Fprintf(os.Stderr, "       cat <files> | %s (one float/integer per line)\n", cmdname)
		os.Exit(1)
	}
	var err error
	var lines []string
	if len(os.Args) > 1 {
		for _, f := range os.Args[1:] {
			ls, err := hio.ReadFileToLines(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read lines from file %s : %v\n", f, err)
				os.Exit(1)
			}
			lines = append(lines, ls...)
		}
	} else {
		lines, err = hio.ReadStdinToLines()
		if err != nil {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read lines from stdin : %v\n", err)
				os.Exit(1)
			}
		}
	}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.Contains(line, "\t") {
			fmt.Printf("Skipping %s\n", line)
			continue
		}
		if strings.Contains(line, "#") {
			fmt.Printf("Skipping %s\n", line)
			continue
		}
		line = strings.TrimSpace(strings.Replace(line, ",", ".", -1))
		line = strings.Replace(line, " ", "", -1)
		asNum, err := strconv.ParseFloat(line, 64)
		n++
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't parse number from %s : %v\n", line, err)
			os.Exit(1)
		}
		sum = sum + asNum
	}
	mean := sum / float64(n)
	fmt.Printf("items   %15d\n", n)
	fmt.Printf("sum     %15.2f\n", sum)
	fmt.Printf("mean    %15.2f\n", mean)
}
