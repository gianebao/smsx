package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gianebao/smsx/app"
)

func init() {
	app.Init()
}

func main() {
	flag.Parse()

	http.HandleFunc("/", app.MessageHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
