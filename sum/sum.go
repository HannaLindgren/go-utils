package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	sum := 0.0
	n := 0
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <files> (one float/integer per line)\n", cmdname)
		os.Exit(1)
	}
	for i := 1; i < len(os.Args); i++ {
		file := os.Args[i]
		bts, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Couldn't read file %s : %v", file, err)
		}
		lines := strings.Split(string(bts), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			if strings.Contains(line, "\t") {
				fmt.Println("Skipping %s", line)
				continue
			}
			line = strings.TrimSpace(strings.Replace(line, ",", ".", -1))
			line = strings.Replace(line, " ", "", -1)
			asNum, err := strconv.ParseFloat(line, 64)
			n++
			if err != nil {
				log.Fatalf("Couldn't parse number from %s : %v", line, err)
			}
			sum = sum + asNum
		}
	}
	mean := sum / float64(n)
	fmt.Printf("items   %15d\n", n)
	fmt.Printf("sum     %15.2f\n", sum)
	fmt.Printf("mean    %15.2f\n", mean)
}
