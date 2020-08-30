package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("https://www.google.com/robots.txt")
	if err != nil {
		log.Panicln(err)
	}
	// Print HTTP status code
	fmt.Println(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()
	fmt.Println(string(body))
}
