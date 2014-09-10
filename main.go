package main

import (
	//"bitbucket.org/alexisjanvier/deployedpr/deployment"
	"encoding/json"
	"fmt"
	"net/http"
)

type Deployment struct {
	GitRef string
	Branch string
	Tag    string
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var testJson string
		testJson = `{"GitRef":"alexisjanvier/projettest","Branch":"master", "Tag":""}`
		byt := []byte(testJson)
		var newDeploy Deployment

		if err := json.Unmarshal(byt, &newDeploy); err != nil {
			panic(err)
		}
		fmt.Println(newDeploy)

		fmt.Fprintf(w, "Welcome to the Deployed Pull Requests Webservice (you asked %q)", r.URL.Path[1:])
	} else {
		http.Error(w, "You must send your request in POST.", 405)
	}
}

func main() {
	http.HandleFunc("/", processRequest)
	http.ListenAndServe(":8080", nil)
}
