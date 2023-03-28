package strings

import (
	"fmt"
	"regexp"
)

func ExampleTokenizer1() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split("") {
		fmt.Printf("<%s>\n", w)
	}
	// Output:
}

func ExampleTokenizer2() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split(" ") {
		fmt.Printf("<%s>\n", w)
	}
	// Output:
	// < >
}

func ExampleTokenizer3() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split("hej du") {
		fmt.Printf("<%s>\n", w)
	}
	// Output:
	// <hej>
	// < >
	// <du>
}

func ExampleTokenizer4() {
	tk := RegexpTokenizer{delimRE: regexp.MustCompile(`[ .,/()&#!?]+`)}
	for _, w := range tk.Split(" -s") {
		fmt.Printf("<%s>\n", w)
	}
	// Output:
	// < >
	// <-s>
}

func ExampleTokenizer5() {
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

func ExampleTokenizer6() {
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
