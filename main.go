package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintf(w, "Welcome to the Deployed Pull Requests Webservice (you asked %q)", r.URL.Path[1:])
	} else {
		http.Error(w, "You must send your request in POST.", 405)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
