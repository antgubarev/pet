package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/antgubarev/gorelex/internal"
)

func main() {
	webHookUri := fmt.Sprintf("/webhook-%s", os.Getenv("PET_WEBHOOK_TOKEN"))
	bot := internal.NewBot(os.Getenv("PET_BOT_TOKEN"))
	if err := bot.InitWebHook(os.Getenv("PET_DOMAIN") + webHookUri); err != nil {
		log.Fatal(err)
	}
	if err := bot.InitCmds(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/count", func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading body: %v", err)
			return
		}
		count := internal.CountWordsInString(string(body))
		rw.Write([]byte(strconv.Itoa(count)))
	})

	http.HandleFunc(webHookUri, func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("WebHook parse: %v", err)
			return
		}
		var update internal.Update
		if err := json.Unmarshal(body, &update); err != nil {
			log.Printf("WebHook update Unmarshal: %v", err)
			return
		}
		bot.WebHook(&update)
	})

	listen := flag.String("port", "8080", "listen port")
	flag.Parse()
	fmt.Println("start the server")
	if err := http.ListenAndServe(":"+*listen, nil); err != nil {
		panic(err)
	}
}
