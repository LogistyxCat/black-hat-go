package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type status struct {
	Message string
	Status  string
}

func main() {
	res, err := http.Post(
		"http://127.0.0.1:12345/ping",
		"application/json",
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	var s status
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	log.Printf("%s -> %s\n", s.Status, s.Message)
}
