package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

func main() {
	cmdname := filepath.Base(os.Args[0])
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		os.Exit(1)
	}

	compared := make(map[string]bool)

	for _, f1 := range os.Args[1:] {
		for _, f2 := range os.Args[1:] {

			if f1 == f2 {
				//fmt.Printf("Files %s and %s are the same file\n", f1, f2)
				continue
			}

			tmp := []string{f1, f2}
			sort.Slice(tmp, func(i, j int) bool { return tmp[i] < tmp[j] })
			compID := strings.Join(tmp, ", ")

			if _, ok := compared[compID]; !ok {

				compared[compID] = true

				bts1, err := ioutil.ReadFile(f1)
				if err != nil {
					log.Fatalf("Couldn't load file %s : %v", f1, err)
				}
				bts2, err := ioutil.ReadFile(f2)
				if err != nil {
					log.Fatalf("Couldn't load file %s : %v", f2, err)
				}
				if reflect.DeepEqual(bts1, bts2) {
					fmt.Printf("Files %s and %s are identical\n", f1, f2)
				} else {
					fmt.Printf("Files %s and %s differ\n", f1, f2)
				}
			} else {
				//fmt.Printf("Files %s and %s already compared\n", f1, f2)
			}
		}
	}
}
