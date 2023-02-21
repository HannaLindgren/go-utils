package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

const cmdname = "rename_files"

func exists(f string) bool {
	_, err := os.Stat(f)
	return !os.IsNotExist(err)
}

func main() {
	//prompt := flag.Bool("p", false, "Prompt before overwriting existing files")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Rename multiple files using regexp")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <from regexp> <to string> <files>\n", cmdname)
		// fmt.Fprintln(os.Stderr, "\nOptional flags:")
		// flag.PrintDefaults()
	}

	flag.Usage = printUsage

	flag.Parse()

	if flag.NArg() < 3 {
		printUsage()
		os.Exit(1)
	}

	fromS := flag.Args()[0]
	toS := flag.Args()[1]
	files := flag.Args()[2:]

	fromRE, err := regexp.Compile(fromS)
	if err != nil {
		log.Fatalf("Regexp compile failed: %v", err)
	}

	for _, oldF := range files {
		newF := fromRE.ReplaceAllString(oldF, toS)
		if newF == oldF {
			fmt.Fprintln(os.Stderr, "Skipping", oldF)
			continue
		}
		fmt.Fprintln(os.Stdout, oldF, "=>", newF)
		err := os.Rename(oldF, newF)
		if err != nil {
			log.Fatalf("Rename failed: %v", err)
		}
	}

}
