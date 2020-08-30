package main

/*
TO DO:
	[x] Limt regex to files only, removing directories
	[x] Locate and flag recently accessed and/or modified files
	[ ] Make var recentDate take user input for preferred date range
*/

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// Calculate recentDate once so it doesn't take unnecessary computations
// Default value is within the last 14 days (2 weeks)
var recentDate = time.Now().AddDate(0, 0, -14)

// This list of regexes can be expanded
var regexes = []*regexp.Regexp{
	regexp.MustCompile(`(?i)user`),
	regexp.MustCompile(`(?i)password`),
	regexp.MustCompile(`(?i)kdb`),
	regexp.MustCompile(`(?i)login`),

	// New, custom regexes here
	regexp.MustCompile(`(?i)rsa`),
	regexp.MustCompile(`(?i)key`),
	regexp.MustCompile(`(?i)psafe`),
}

func walkFn(path string, f os.FileInfo, err error) error {
	// We aren't interested in directory names at this time
	if f.IsDir() {
		return nil
	}

	for _, r := range regexes {
		if r.MatchString(path) {
			fmt.Printf("[+] HIT: %s\n", path)
		}
	}

	// Show recently accessed files
	if f.ModTime().After(recentDate) {
		fmt.Printf("[*] Recently modified: %s\n", path)
	}

	return nil
}

func main() {
	var root string
	if len(os.Args) != 1 {
		root = os.Args[1]
	} else {
		root = "/"
	}

	if err := filepath.Walk(root, walkFn); err != nil {
		log.Panicln(err)
	}
}
