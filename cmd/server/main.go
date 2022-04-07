package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/antgubarev/gorelex/internal"
)

func main() {
	http.HandleFunc("/count", func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading body: %v", err)
			return
		}
		count := internal.CountWordsInString(string(body))
		rw.Write([]byte(strconv.Itoa(count)))
	})

	fmt.Println("server start")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		panic(err)
	}
}
