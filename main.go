package main

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
)

const (
	PROTOCOL_VERSION int = 1
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

func serveLib(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/lib" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	t := template.Must(template.ParseFiles("lib.min.js"))

	data := map[string]string{
		"version":  strconv.Itoa(PROTOCOL_VERSION),
		"endpoint": r.Host,
	}

	t.Execute(w, data)
}

func serveEvents(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/event" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.Write([]byte("{\"time\":1}"))

}

func main() {

	//load settings from ini and params
	configSetup()

	go h.run()

	http.HandleFunc("/test", serveHome)
	http.HandleFunc("/lib", serveLib)
	http.HandleFunc("/event", serveEvents)
	http.HandleFunc("/", serveWs)

	log.Print("Listen on ", config.BindAddress)
	err := http.ListenAndServe(config.BindAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
