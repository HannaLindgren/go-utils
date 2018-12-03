package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/HannaLindgren/go-utils/scripts/util"
)

func readFromFilesOrStdin(files []string) (map[string]int64, int, error) {
	freq := make(map[string]int64)
	tot := 0
	if len(files) > 0 {
		for _, f := range files {
			r, fh, err := util.GetFileReader(f)
			defer fh.Close()
			if err != nil {
				return freq, 0, err
			}
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				s := scanner.Text()
				freq[s]++
				tot++
			}

		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			s := scanner.Text()
			freq[s]++
			tot++
		}
	}
	return freq, tot, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func sortByValue(freqMap map[string]int64) []string {
	res := []string{}
	for s := range freqMap {
		if !contains(res, s) {
			res = append(res, s)
		}
	}
	sort.Slice(res, func(i, j int) bool { return freqMap[res[i]] > freqMap[res[j]] })
	return res
}

func printHelpAndExit() {
	cmdname := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s <files> or cat <files> | %s\n", cmdname, cmdname)
	fmt.Fprintf(os.Stderr, "Switches:\n")
	fmt.Fprintf(os.Stderr, " -h print help and exit\n")
	fmt.Fprintf(os.Stderr, " -r print frequency on the right hand side (default: false)\n")
	fmt.Fprintf(os.Stderr, " -p print percentage (default: false)\n")
	os.Exit(1)
}

func prcntFmt(value int64, total int) string {
	return fmt.Sprintf("%.2f%%", math.Round(float64(value*100))/float64(total))
}

func main() {
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-h") {
		printHelpAndExit()
	}
	args := os.Args[1:]
	freqRight := false
	percentage := false
	if len(args) > 0 {
		if strings.HasPrefix(args[0], "-r") {
			freqRight = true
			args = args[1:]
		} else if strings.HasPrefix(args[0], "-p") {
			percentage = true
			args = args[1:]
		} else {
			printHelpAndExit()
		}

	}

	freq, total, err := readFromFilesOrStdin(args)
	if err != nil {
		log.Fatalf("Couldn't compute : %v", err)
	}
	for _, s := range sortByValue(freq) {
		if freqRight {
			if percentage {
				fmt.Printf("%v\t%v\t%s\n", s, freq[s], prcntFmt(freq[s], total))
			} else {
				fmt.Printf("%v\t%v\n", s, freq[s])
			}
		} else if percentage {
			fmt.Printf("%v\t%s\t%v\n", freq[s], prcntFmt(freq[s], total), s)
		} else {
			fmt.Printf("%v\t%v\n", s, freq[s])
		}
	}
}
