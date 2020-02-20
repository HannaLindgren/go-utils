package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func isFile(fName string) bool {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return false
	}
	return true
}

func readFile(fName string) []string {
	b, err := ioutil.ReadFile(filepath.Clean(fName))
	if err != nil {
		log.Fatalf("%v", err)
		return []string{}
	}
	s := strings.TrimSuffix(string(b), "\n")
	s = strings.Replace(s, "\r", "", -1)
	return strings.Split(s, "\n")
}

type pair struct {
	s1 string
	s2 string
}

type repair struct {
	from *regexp.Regexp
	to   string
}

// https://en.wikipedia.org/wiki/Romanization_of_Greek#Modern_Greek

var mapRegexps = []repair{
	repair{from: regexp.MustCompile(`(^|[\s/()'".!?-])Γ[Κκ](.+)`), to: "${1}G${2}"},
	repair{from: regexp.MustCompile(`(^|[\s/()'".!?-])Μ[Ππ](.+)`), to: "${1}B${2}"},
	repair{from: regexp.MustCompile(`(^|[\s/()'".!?-])Ν[Ττ](.+)`), to: "${1}D${2}"},
	repair{from: regexp.MustCompile(`(^|[\s/()'".!?-])(?i)γκ(.+)`), to: "${1}g${2}"},
	repair{from: regexp.MustCompile(`(^|[\s/()'".!?-])(?i)μπ(.+)`), to: "${1}b${2}"},
	repair{from: regexp.MustCompile(`(^|[\s/()'".!?-])(?i)ντ(.+)`), to: "${1}d${2}"},
}

var maptable = []pair{
	pair{s1: "αι", s2: "ai"},
	pair{s1: "ει", s2: "ei"},
	pair{s1: "οι", s2: "oi"},
	pair{s1: "υι", s2: "yi"},

	pair{s1: "αυ", s2: "au"},
	pair{s1: "ευ", s2: "eu"},
	pair{s1: "ου", s2: "ou"},

	pair{s1: "αύ", s2: "au"},
	pair{s1: "εύ", s2: "eú"},
	pair{s1: "ού", s2: "oú"},
	pair{s1: "άυ", s2: "áu"},
	pair{s1: "έυ", s2: "éu"},
	pair{s1: "όυ", s2: "óu"},

	pair{s1: "ήυ", s2: "íy"},
	pair{s1: "υί", s2: "yí"},
	pair{s1: "ηυ", s2: "iy"},

	pair{s1: "ωυ", s2: "oy"},
	pair{s1: "ώυ", s2: "óy"},

	pair{s1: "μμπ", s2: "mb"},
	pair{s1: "νντ", s2: "nd"},

	pair{s1: "ά", s2: "á"},
	pair{s1: "έ", s2: "é"},
	pair{s1: "ή", s2: "í"},
	pair{s1: "ί", s2: "í"},
	pair{s1: "ύ", s2: "í"},
	pair{s1: "ό", s2: "ó"},
	pair{s1: "ώ", s2: "ó"},
	pair{s1: "ϊ", s2: "ï"},
	pair{s1: "ΐ", s2: "ḯ"},
	pair{s1: "ϋ", s2: "ü"},
	pair{s1: "ΰ", s2: "ǘ"},

	pair{s1: "α", s2: "a"},
	pair{s1: "β", s2: "v"},
	pair{s1: "γγ", s2: "ng"},
	pair{s1: "γκ", s2: "nk"},
	pair{s1: "γξ", s2: "nx"},
	pair{s1: "γχ", s2: "nch"},
	pair{s1: "γ", s2: "g"},
	pair{s1: "δ", s2: "d"},
	pair{s1: "ε", s2: "e"},
	pair{s1: "ζ", s2: "z"},

	pair{s1: "η", s2: "i"},
	pair{s1: "θ", s2: "th"},
	pair{s1: "ι", s2: "i"},
	pair{s1: "κ", s2: "k"},
	pair{s1: "λ", s2: "l"},
	pair{s1: "μ", s2: "m"},
	pair{s1: "ν", s2: "n"},
	pair{s1: "ξ", s2: "x"},
	pair{s1: "ο", s2: "o"},
	pair{s1: "π", s2: "p"},
	pair{s1: "ρ", s2: "r"},
	pair{s1: "ς", s2: "s"},
	pair{s1: "σ", s2: "s"},
	pair{s1: "τ", s2: "t"},
	pair{s1: "υ", s2: "y"},
	pair{s1: "φ", s2: "f"},
	pair{s1: "χ", s2: "ch"},
	pair{s1: "ψ", s2: "ps"},
	pair{s1: "ω", s2: "o"},
}

var commonCharsRE = regexp.MustCompile("[A-Za-z0-9()@΄$ï*_]")

var commonChars = map[string]bool{
	" ":  true,
	"\t": true,
	",":  true,
	".":  true,
	"?":  true,
	"!":  true,
	"–":  true,
	"-":  true,
	":":  true,
	";":  true,
	"&":  true,
	"/":  true,
	"'":  true,
}

func convert(s string) (string, error) {
	for _, re := range mapRegexps {
		s = re.from.ReplaceAllString(s, re.to)
	}

	intAll := []pair{}
	for _, p := range maptable {
		intAll = append(intAll, p)
		intAll = append(intAll, pair{s1: upcaseInitial(p.s1), s2: upcaseInitial(p.s2)})
		intAll = append(intAll, pair{s1: upcase(p.s1), s2: upcase(p.s2)})
		if len([]rune(p.s1)) > 2 {
			intAll = append(intAll, pair{s1: upcaseTwoInitials(p.s1), s2: upcaseInitial(p.s2)})
		}
	}
	res, err := innerConvert(intAll, s, true)
	if err != nil {
		return "", err
	}
	return res, nil
}

func innerConvert(chsAll []pair, s string, requireAllMapped bool) (string, error) {
	sOrig := s
	res := []string{}
	for len(s) > 0 {
		sStart := s
		head := string([]rune(s)[0])
		for _, p := range chsAll {
			if strings.HasPrefix(s, p.s1) {
				res = append(res, p.s2)
				s = strings.TrimPrefix(s, p.s1)
				break
			}
			// check for common chars
			if _, ok := commonChars[head]; ok {
				res = append(res, head)
				s = strings.TrimPrefix(s, head)
				break
			}
			if commonCharsRE.MatchString(head) {
				res = append(res, head)
				s = strings.TrimPrefix(s, head)
				break
			}
		}
		if s == sStart { // nothing found to map for this prefix
			if requireAllMapped {
				return "", fmt.Errorf("Couldn't convert '%s' in '%s'", s, sOrig)
			} else if s == sStart {
				res = append(res, head)
				s = strings.TrimPrefix(s, head)
			}
		}
	}
	return strings.Join(res, ""), nil
}

func upcase(s string) string {
	return strings.ToUpper(s)
}

func upcaseInitial(s string) string {
	runes := []rune(s)
	head := ""
	if len(runes) > 0 {
		head = strings.ToUpper(string(runes[0]))
	}
	tail := ""
	if len(runes) > 1 {
		tail = strings.ToLower(string(runes[1:]))
	}
	//fmt.Println("??? 2", len(runes), s, head, tail)
	return head + tail
}

func upcaseTwoInitials(s string) string {
	runes := []rune(s)
	head := ""
	if len(runes) > 0 {
		head = strings.ToUpper(string(runes[0:2]))
	}
	tail := ""
	if len(runes) > 2 {
		tail = strings.ToLower(string(runes[2:]))
	}
	//fmt.Println("??? 3", len(runes), s, head, tail)
	return head + tail
}

func process(s string) {
	res, err := convert(s)
	if err != nil {
		if *failOnError {
			log.Fatalf(fmt.Sprintf("%v", err))
		} else {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("ERROR %s\t%v\n", s, err))
			return
		}
	}
	if *echoInput {
		fmt.Printf("%s\t%s\n", s, res)
	} else {
		fmt.Printf("%s\n", res)
	}
}

var echoInput, failOnError *bool

func main() {

	cmdname := filepath.Base(os.Args[0])
	echoInput = flag.Bool("e", false, "Echo input (default: false)")
	failOnError = flag.Bool("f", false, "Fail on error (default: false)")
	help := flag.Bool("h", false, "Print help and exit")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Transliteration from Greek to Latin script.")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, cmdname+" <input file(s)>")
		fmt.Fprintln(os.Stderr, cmdname+" <input string(s)>")
		fmt.Fprintln(os.Stderr, "cat <input file(s)> | "+cmdname)
		fmt.Fprintln(os.Stderr, "\nOptional flags:")
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}

	flag.Parse()

	if *help { // if flag.NArg() < 1 {
		printUsage()
		os.Exit(0)
	}

	if len(flag.Args()) > 0 {
		for _, arg := range flag.Args() {
			if isFile(arg) {
				for _, line := range readFile(arg) {
					process(line)
				}
			} else {
				process(arg)
			}
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			s := scanner.Text()
			process(s)
		}
	}
}
