package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Post("http://127.0.0.1:8080/count", "application/text", bytes.NewBuffer([]byte("string with four words")))
	if err != nil {
		log.Fatalf("send requset to count: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("parse response body: %v", err)
	}

	log.Print("count: ", string(body))
}
