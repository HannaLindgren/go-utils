package strings

import (
	"fmt"
	//"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/runenames"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func NFKC(s string) string {
	return string(norm.NFKC.Bytes([]byte(s)))
}

func NFKD(s string) string {
	return string(norm.NFKD.Bytes([]byte(s)))
}

func NFC(s string) string {
	return string(norm.NFC.Bytes([]byte(s)))
}

func NFD(s string) string {
	return string(norm.NFD.Bytes([]byte(s)))
}

// UnicodeForR Returns unicode for the input rule
func UnicodeForR(r rune) string {
	uc := fmt.Sprintf("%U", r)
	return fmt.Sprintf("\\u%s", uc[2:])
}

const newline rune = '\n'

// UnicodeFor Returns a list of unicodes for each input rune
func UnicodeFor(s string) []string {
	res := []string{}
	for _, r := range s {
		res = append(res, UnicodeForR(r))
		if r == newline {
			fmt.Println()
		}
	}
	return res
}

var ucNumberRe = regexp.MustCompile(`^(?:\\u|[uU][+])([a-fA-F0-9]{4})$`)

var hardwiredNames = map[rune]string{
	newline: "NEWLINE",
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

// UnicodeInfo holds a set of unicode-related information for a rune
type UnicodeInfo struct {
	// Rune is the input rune
	Rune rune

	// A string representation of the input rune ('NEWLINE' for newline, 'TAB' for tab, etc)
	String string

	// The unicode number
	Unicode string

	// The character name
	CharName string

	// The codeblock
	CodeBlock string
}

// UnicodeProcessor
type UnicodeProcessor struct {
	nfc                       bool
	nfd                       bool
	convertFromUnicodeNumbers bool
}

func (up *UnicodeProcessor) inhibitSpecialChar(r rune) bool {
	_, ok := hardwiredNames[r]
	return ok
}

func (up *UnicodeProcessor) blockFor(r rune) string {
	for s, t := range unicode.Scripts {
		if unicode.In(r, t) {
			return s
		}
	}
	return "<UNDEF>"
}

func (up *UnicodeProcessor) nameFor(r rune) string {
	if name, ok := hardwiredNames[r]; ok {
		return name
	}
	return runenames.Name(r)
}

func (up *UnicodeProcessor) normalize(s string) string {
	if up.nfc {
		return NFC(s)
	} else if up.nfd {
		return NFD(s)
	}
	return s
}

// UnicodeInfoR Returns tab-separated unicode information for each input rune
func (up *UnicodeProcessor) UnicodeInfoR(r rune) UnicodeInfo {
	name := up.nameFor(r)
	uc := UnicodeForR(r)
	block := up.blockFor(r)
	s := string(r)
	if up.inhibitSpecialChar(r) {
		s = ""
	}
	//return fmt.Sprintf("%s\t%s\t%s\t%s\n", thisS, uc, name, block)
	return UnicodeInfo{
		Rune:      r,
		String:    s,
		Unicode:   uc,
		CharName:  name,
		CodeBlock: block,
	}
}

// UnicodeInfo Creates a list of tab-separated unicode information for each input rune
func (up *UnicodeProcessor) UnicodeInfo(s string) []UnicodeInfo {
	res := []UnicodeInfo{}
	sx := s
	if up.convertFromUnicodeNumbers {
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
	for _, r := range up.normalize(sx) {
		res = append(res, up.UnicodeInfoR(r))
	}
	return res
}
