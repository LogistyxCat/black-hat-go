package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

// GLOBAL VARIABLES
var addr = ":8000"

type router struct{}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/a":
		fmt.Fprintf(w, "Executing /a\n")
	case "/b":
		fmt.Fprintf(w, "Executing /b\n")
	case "/c":
		fmt.Fprintf(w, "Executing /c\n")
	default:
		http.Error(w, "404 Not Found", 404)
	}
}

func init() {
	help := flag.Bool("h", false, "Print usage and exit")
	flag.StringVar(&addr, "addr", addr, "Address to Listen on")
	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
	}
}

// usage prints command line options
func usage() {
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	var r router
	log.Fatal(http.ListenAndServe(addr, &r))
}
