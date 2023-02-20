package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HannaLindgren/go-utils/io"
	"github.com/HannaLindgren/go-utils/unicode"
)

type TokInput struct {
	Input      string          `json:"input"`
	TokenCount int             `json:"count"`
	Tokens     []unicode.Token `json:"tokens"`
}

func main() {
	cmdname := filepath.Base(os.Args[0])
	nfc := flag.Bool("nfc", false, "NFC -- Canonical composition on all input (default false)")
	nfd := flag.Bool("nfd", false, "NFD -- Canonical decomposition on all input (default false)")
	outFmt := flag.String("o", "t", "Output type [t=tab-separated (simplified oneline format); j=json]")
	splitInputLines := flag.Bool("l", false, "Split input by newline before tokenizing (default for tab-separated output)")
	skipWhiteSpace := flag.Bool("sw", false, "Skip white space (default for tab-separated output)")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Utility script to tokenise strings based on their unicode block.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "Usage: %s <files>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       %s <string>\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       cat <file> | %s\n", cmdname)
		fmt.Fprintf(os.Stderr, "       or\n")
		fmt.Fprintf(os.Stderr, "       echo <string> | %s\n", cmdname)
		// fmt.Fprintln(os.Stderr, cmdname+" <flags> <file1> <file2>")
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	if *nfd && *nfc {
		fmt.Fprintf(os.Stderr, "nfc and nfd options cannot be combined\n")
		printUsage()
		os.Exit(0)
	}

	var jsonO, tabO bool
	switch *outFmt {
	case "j":
		jsonO = true
	case "t":
		tabO = true
		t := true
		skipWhiteSpace = &t
		splitInputLines = &t
	default:
		fmt.Fprintf(os.Stderr, "Invalid output format: %v\n", *outFmt)
		printUsage()
		os.Exit(1)

	}

	toker := unicode.Tokenizer{
		UP: unicode.Processor{
			NFC: *nfc,
			NFD: *nfd,
		},
		SkipWhiteSpace: *skipWhiteSpace,
	}

	var process = func(s0 string) {
		input := []string{}
		if *splitInputLines {
			input = strings.Split(s0, "\n")
		} else {
			input = append(input, s0)
		}
		jsonAllRes := []TokInput{}
		for _, s := range input {
			tokens := toker.Tokenize(s)
			jsonRes := TokInput{Input: s, TokenCount: len(tokens)}
			tabRes := []string{s, strconv.Itoa(len(tokens))}
			for _, t := range tokens {
				if jsonO {
					jsonRes.Tokens = append(jsonRes.Tokens, t)
					//fmt.Printf("<token type='%s'>%s</token>\n", t.UnicodeBlock, t.String)
				} else if tabO {
					tabRes = append(tabRes, t.String)
				}
			}
			if tabO {
				fmt.Println(strings.Join(tabRes, "\t"))
			} else if jsonO {
				jsonAllRes = append(jsonAllRes, jsonRes)
			}
		}
		if jsonO {
			bts, err := json.MarshalIndent(jsonAllRes, " ", " ")
			if err != nil {
				log.Fatalf("%v", err)
			}
			fmt.Println(string(bts))
		}
	}

	if len(flag.Args()) > 0 {
		for _, arg := range os.Args[1:] {
			if io.IsFile(arg) {
				text, err := io.ReadFileToString(arg)
				if err != nil {
					log.Fatalf("%v", err)
				}
				process(text)
			} else {
				process(arg)
			}
		}
	} else {
		text, err := io.ReadStdinToString()
		if err != nil {
			log.Fatalf("%v", err)
		}
		process(text)
	}
}
