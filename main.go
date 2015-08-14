package main

import (
	"log"
	"net/http"
	"text/template"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/test" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles("ws.html")).Execute(w, r.Host)
}

func main() {

	//load settings from ini and params
	configSetup()

	go h.run()

	http.HandleFunc("/test", serveHome)
	http.HandleFunc("/", serveWs)

	log.Print("Listen on ", config.BindAddress)
	err := http.ListenAndServe(config.BindAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
