package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	hio "github.com/HannaLindgren/go-utils/io"
)

const cmdname = "sum"

var nItems int
var sum float64

func process(lines []string) {
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 || strings.HasPrefix(l, "#") {
			continue
		}
		l = strings.Replace(l, ",", ".", -1)
		nItems++
		f, err := strconv.ParseFloat(l, 10)
		if err != nil {
			log.Fatalf("Parse float failed for %v: %v", l, err)
		}
		sum += f
	}
}

func main() {
	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Sum of numeric input")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <file> or cat <file> | %s\n", cmdname)
	}

	if strings.HasPrefix(os.Args[0], "-h") {
		printUsage()
		os.Exit(0)
	}

	if len(os.Args) == 1 {
		lines, err := hio.ReadStdinToLines()
		if err != nil {
			log.Fatalf("Read failed for stdin: %v", err)
		}
		process(lines)
	} else {
		for _, f := range os.Args[1:] {
			lines, err := hio.ReadFileToLines(f)
			if err != nil {
				log.Fatalf("Read failed for file %s: %v", f, err)
			}
			process(lines)
		}
	}

	mean := sum / float64(nItems)
	fmt.Fprintf(os.Stdout, "sum:   %v\n", sum)
	fmt.Fprintf(os.Stdout, "items: %v\n", nItems)
	fmt.Fprintf(os.Stdout, "mean:  %v\n", mean)

}
