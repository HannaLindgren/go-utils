package main

// https://en.wikipedia.org/wiki/Romanization_of_Russian

import (
	"fmt"
	"log"
	"os"
	"strings"
	"io/ioutil"
)

func IsFile(fName string) bool {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadFile(fName string) []string {
	b, err := ioutil.ReadFile(fName)
	if err != nil {
		log.Fatalf("%v", err)
		return []string{}
	}
	return strings.Split(strings.TrimSuffix(string(b), "\n"), "\n")
}


type pair struct {
	s1 string
	s2 string
}

var charset = []pair{
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
	pair{s1: "х", s2: "ch"},
	pair{s1: "ц", s2: "ts"},
	pair{s1: "ч", s2: "tj"},
	pair{s1: "ш", s2: "sj"},
	pair{s1: "щ", s2: "sjtj"},
	pair{s1: "ъ", s2: ""},
	pair{s1: "ь", s2: ""},
	pair{s1: "э", s2: "e"},
	pair{s1: "ю", s2: "ju"},
	pair{s1: "я", s2: "ja"},
	pair{s1: "ы", s2: "y"},
	
	// common
	pair{s1: " ", s2: " "},
	pair{s1: ",", s2: ","},
	pair{s1: ".", s2: "."},
	pair{s1: "?", s2: "?"},
	pair{s1: "!", s2: "!"},
}

func convert(chsAll []pair, s string) string {
	sOrig := s
	res := []string{}
	for len(s) > 0 {
		sStart := s
		for _, p := range chsAll {
			if strings.HasPrefix(s, p.s1) {
				//fmt.Println(p.s1, "->", p.s2)
				res = append(res, p.s2)
				s = strings.TrimPrefix(s, p.s1)
				break
			}
		}
		if s == sStart {
			log.Fatalf("Couldn't convert '%s' at '%s'", sOrig, s)
		}
	}
	return strings.Join(res, "")
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

func main() {

	chsAll := []pair{}
	for _, p := range charset {
		chsAll = append(chsAll, p)
		chsAll = append(chsAll, pair{s1: upcaseInitial(p.s1), s2: upcaseInitial(p.s2)})
	}
	args := os.Args[1:]
	//args := []string{"Антон Чехов"}
	for _, arg := range args {
		if IsFile(arg) {
			for _, line := range ReadFile(arg) {
				res := convert(chsAll, line)
				//fmt.Printf("%s\t%s\n", line, res)
				fmt.Printf("%s\n", res)	
			}	
		} else { 
			res := convert(chsAll, arg)
			//fmt.Printf("%s\t%s\n", arg, res)
			fmt.Printf("%s\n", res)
		}
	}
}
