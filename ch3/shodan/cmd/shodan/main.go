package main

import (
	"bhg-prsnl/ch3/shodan/shodan"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: shodan searchterm")
	}

	apiKey := os.Getenv("SHODAN_API_KEY")
	s := shodan.New(apiKey)
	info, err := s.APIInfo()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf(
		"Query Credits: %d\nScan Credits: %d\n\n",
		info.QueryCredits,
		info.ScanCredits,
	)

	prof, err := s.GetProfileInfo()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf(
		"Display Name: %s\nMember: %t\nExport Credits: %d\nCreation Date: %s\n\n",
		prof.DisplayName,
		prof.Member,
		prof.ExportCredits,
		prof.CreatedDate,
	)
}
