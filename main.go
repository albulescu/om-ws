package main

import (
	"log"
	"net/http"
)

func main() {

	//load settings from ini and params
	configSetup()

	go h.run()

	http.HandleFunc("/", serveWs)
	log.Print("Listen on ", config.BindAddress)
	err := http.ListenAndServe(config.BindAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
