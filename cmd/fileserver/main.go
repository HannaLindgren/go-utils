package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", "", "the directory of static file to host")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage: go run main.go <port> <directory>\n")
		fmt.Fprint(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *directory == "" {
		fmt.Fprint(os.Stderr, "Missing required flag -d (directory)\n")
		flag.Usage()
		os.Exit(1)
	}

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
