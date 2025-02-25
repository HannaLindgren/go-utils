package strings

import (
	"fmt"
	"regexp"
)

func ExampleRegexpTokenizer_t1() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split("") {
		fmt.Printf("<%s>\n", w)
	}
	// Output:
}

func ExampleRegexpTokenizer_t2() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split(" ") {
		fmt.Printf("<%s>\n", w)
	}
	// Output:
	// < >
}

func ExampleRegexpTokenizer_t3() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split("hej du") {
		fmt.Printf("<%s>\n", w)
	}
	// Output:
	// <hej>
	// < >
	// <du>
}

func ExampleRegexpTokenizer_t4() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split(" -s") {
		fmt.Printf("<%s>\n", w)
	}
	// Output:
	// < >
	// <-s>
}

func ExampleRegexpTokenizer_t5() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split("jag-är!& -en liten apa") {
		fmt.Printf("<%s>\n", w)
		//fmt.Fprintf(os.Stderr, "<%s>\n", w)
	}
	// Output:
	// <jag-är>
	// <!& >
	// <-en>
	// < >
	// <liten>
	// < >
	// <apa>
}

func ExampleRegexpTokenizer_t6() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split("jag-är!& -en liten apa!") {
		fmt.Printf("<%s>\n", w)
		//fmt.Fprintf(os.Stderr, "<%s>\n", w)
	}
	// Output:
	// <jag-är>
	// <!& >
	// <-en>
	// < >
	// <liten>
	// < >
	// <apa>
	// <!>
}
