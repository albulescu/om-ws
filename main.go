package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":9998", "http service address")

func main() {

	//load settings from ini and params
	configSetup()

	go h.run()

	http.HandleFunc("/ws", serveWs)
	log.Print("Listen on ", config.BindAddress)
	err := http.ListenAndServe(config.BindAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
