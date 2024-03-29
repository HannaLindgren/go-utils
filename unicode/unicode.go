package unicode

import (
	"fmt"
	//"golang.org/x/text/transform"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/runenames"
)

func NFKC(s string) string {
	//norm, _, _ := transform.String(norm.NFKC, s)
	return string(norm.NFKC.Bytes([]byte(s)))
}

func NFKD(s string) string {
	//norm, _, _ := transform.String(norm.NFKD, s)
	// return norm
	return string(norm.NFKD.Bytes([]byte(s)))
}

func NFC(s string) string {
	//norm, _, _ := transform.String(norm.NFC, s)
	// return norm
	return string(norm.NFC.Bytes([]byte(s)))
}

func NFD(s string) string {
	//norm, _, _ := transform.String(norm.NFD, s)
	// return norm
	return string(norm.NFD.Bytes([]byte(s)))
}

// UnicodeForR Returns unicode for the input rule
func UnicodeForR(r rune) string {
	uc := fmt.Sprintf("%U", r)
	return fmt.Sprintf(`\u%s`, uc[2:])
}

const Newline rune = '\n'

// UnicodeFor Returns a list of unicodes for each input rune
func UnicodeFor(s string) []string {
	res := []string{}
	for _, r := range s {
		res = append(res, UnicodeForR(r))
		// if r == Newline {
		// 	fmt.Println()
		// }
	}
	return res
}

var ucNumberRe = regexp.MustCompile(`^(?:\\u|[uU][+])([a-fA-F0-9]{4})$`)

var specialChars = map[rune]string{
	Newline: "NEWLINE",
	'	':     "TAB",
}

func isUcNumber(s string) bool {
	return ucNumberRe.MatchString(s)
}

func ucNumber2String(s string) string {
	is := ucNumberRe.ReplaceAllString(s, "$1")
	i, err := strconv.ParseInt(is, 16, 32)
	if err != nil {
		log.Fatalf("Couldn't parse unicode number from input string '%s' : %v", s, err)
	}
	r := rune(i)
	return string(r)
}

// Info holds a set of unicode-related information for a rune
type Info struct {
	// A string representation of the input rune (for special newline and tab, the string representation is empty in this implementation)
	String string

	// The unicode number
	Unicode string

	// The character name
	CharName string

	// The codeblock
	CodeBlock string
}

// Processor
type Processor struct {
	NFC                       bool
	NFD                       bool
	ConvertFromUnicodeNumbers bool
}

// BlockFor returns the name of the unicode script where the input rune belongs
func BlockFor(r rune) string {
	for s, t := range unicode.Scripts {
		if unicode.In(r, t) {
			return s
		}
	}
	return "<UNDEF>"
}

// NameFor returns the name of the input rune
func NameFor(r rune) string {
	if name, ok := specialChars[r]; ok {
		return name
	}
	return runenames.Name(r)
}

// Normalize according to the NFC/NFD settings in the UnicodeProcessor
func (p *Processor) Normalize(s string) string {
	if p.NFC {
		return NFC(s)
	} else if p.NFD {
		return NFD(s)
	}
	return s
}

// RuneInfo Returns tab-separated unicode information for each input rune
func (p *Processor) RuneInfo(r rune) Info {
	name := NameFor(r)
	uc := UnicodeForR(r)
	block := BlockFor(r)
	char := string(r)
	if _, inhibitChar := specialChars[r]; inhibitChar {
		char = ""
	}
	//return fmt.Sprintf("%s\t%s\t%s\t%s\n", thisS, uc, name, block)
	return Info{
		String:    char,
		Unicode:   uc,
		CharName:  name,
		CodeBlock: block,
	}
}

// Info Creates a list with unicode information for each input rune
func (p *Processor) UnicodeInfo(s string) []Info {
	res := []Info{}
	sx := s
	if p.ConvertFromUnicodeNumbers {
		tmp := []string{}
		for _, chunk := range strings.Split(sx, " ") {
			if isUcNumber(chunk) {
				tmp = append(tmp, ucNumber2String(chunk))
			} else {
				tmp = append(tmp, chunk)
			}
		}
		sx = strings.Join(tmp, " ")
	}
	for _, r := range p.Normalize(sx) {
		res = append(res, p.RuneInfo(r))
	}
	return res
}
