package main

// References:
// https://en.wikipedia.org/wiki/Romanization_of_Russian
// https://tt.se/tt-spraket/ord-och-begrepp/internationellt/andra-sprak/ryska/

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

type rPair struct {
	from *regexp.Regexp
	to   string
}

var international = []pair{ // https://en.wikipedia.org/wiki/Romanization_of_Russian
	pair{s1: "а", s2: "a"},
	pair{s1: "б", s2: "b"},
	pair{s1: "в", s2: "v"},
	pair{s1: "г", s2: "g"},
	pair{s1: "д", s2: "d"},
	pair{s1: "е", s2: "e"},
	pair{s1: "ё", s2: "e"},
	pair{s1: "ж", s2: "zh"},
	pair{s1: "з", s2: "z"},
	pair{s1: "и", s2: "i"},
	pair{s1: "й", s2: "j"},
	pair{s1: "к", s2: "k"},
	pair{s1: "л", s2: "l"},
	pair{s1: "м", s2: "m"},
	pair{s1: "н", s2: "n"},
	pair{s1: "о", s2: "o"},
	pair{s1: "п", s2: "p"},
	pair{s1: "р", s2: "r"},
	pair{s1: "с", s2: "s"},
	pair{s1: "т", s2: "t"},
	pair{s1: "у", s2: "u"},
	pair{s1: "ф", s2: "f"},
	pair{s1: "х", s2: "kh"},
	pair{s1: "ц", s2: "ts"},
	pair{s1: "ч", s2: "ch"},
	pair{s1: "ш", s2: "sh"},
	pair{s1: "щ", s2: "shch"},
	pair{s1: "ъ", s2: "’"},
	pair{s1: "ы", s2: "y"},
	pair{s1: "ь", s2: "’"},
	pair{s1: "э", s2: "e"},
	pair{s1: "ю", s2: "ju"},
	pair{s1: "я", s2: "ja"},
}

var commonChars = map[string]bool{
	" ": true,
	",": true,
	".": true,
	"?": true,
	"!": true,
	"–": true,
	"-": true,
	":": true,
	";": true,
}

// https://tt.se/tt-spraket/ord-och-begrepp/internationellt/andra-sprak/ryska/
var swePairs = []pair{
	pair{s1: "zh", s2: "zj"},
	pair{s1: "kh", s2: "ch"},
	pair{s1: "ch", s2: "tj"},
	pair{s1: "sh", s2: "sj"},
	//pair{s1: "shch", s2: "sjtj"}, // not needed
	pair{s1: "yu", s2: "ju"},
	pair{s1: "ya", s2: "ja"},
	pair{s1: "ye", s2: "je"}, // ?? will this ever happen?
}

var sweREs = []rPair{ // https://tt.se/tt-spraket/ord-och-begrepp/internationellt/andra-sprak/ryska/

	// must be unambiguous for output

	// rPair{from: regexp.MustCompile(`(?:i)ev([\p{P}]|$)`), to: "(j)ev$1"},
	// rPair{from: regexp.MustCompile(`(?:i)ov([\p{P}]|$)`), to: "ov$1"},
	// rPair{from: regexp.MustCompile(`(?:i)ev([\p{P}]|$)`), to: "jov$1"}, // Gorbatjov

	rPair{from: regexp.MustCompile(`(?i)ky(\b|$)`), to: "kij$2"},
	rPair{from: regexp.MustCompile(`(?i)gy(\b|$)`), to: "gij$2"},
	rPair{from: regexp.MustCompile(`(?i)ay(\b|$)`), to: "aj$2"},
	rPair{from: regexp.MustCompile(`(?i)ey(\b|$)`), to: "ej$2"},
	rPair{from: regexp.MustCompile(`(?i)y(\b|$)`), to: "yj$2"},
}

func convert(s string) (string, error) {
	intAll := []pair{}
	for _, p := range international {
		intAll = append(intAll, p)
		intAll = append(intAll, pair{s1: upcaseInitial(p.s1), s2: upcaseInitial(p.s2)})
		intAll = append(intAll, pair{s1: upcase(p.s1), s2: upcase(p.s2)})
	}
	res, err := innerConvert(intAll, s, true)
	if err != nil {
		return "", err
	}
	if *swedishOutput {
		sweAll := []pair{}
		for _, p := range swePairs {
			sweAll = append(sweAll, p)
			sweAll = append(sweAll, pair{s1: upcaseInitial(p.s1), s2: upcaseInitial(p.s2)})
			sweAll = append(sweAll, pair{s1: upcase(p.s1), s2: upcase(p.s2)})
		}
		res, err := innerConvert(sweAll, res, false)
		if err != nil {
			return "", err
		}
		for _, p := range sweREs {
			res = p.from.ReplaceAllString(res, p.to)
		}
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
	if len(runes) > 0 {
		tail = strings.ToLower(string(runes[1:]))
	}
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

var swedishOutput, echoInput, failOnError *bool

func main() {

	cmdname := filepath.Base(os.Args[0])
	swedishOutput = flag.Bool("s", false, "Swedish (TT style) output (default: international output)")
	echoInput = flag.Bool("e", false, "Echo input (default: false)")
	failOnError = flag.Bool("f", false, "Fail on error (default: false)")
	help := flag.Bool("h", false, "Print help and exit")

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, "Transliteration from Cyrillic to Latin script.")
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
