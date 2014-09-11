package main

import (
	"bitbucket.org/alexisjanvier/deployedpr/deployment"
	"encoding/json"
	"fmt"
	"net/http"
)

func processRequest(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			errorMsg := fmt.Sprintf("%q", e)
			http.Error(w, errorMsg, 400)
		}
	}()
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var newDeploy deployment.Deployment
		if err := decoder.Decode(&newDeploy); err != nil {
			panic(err)
		}
		if jsonValid, erroMsg := newDeploy.IsValid(); !jsonValid {
			panic(erroMsg)
		}
		project, err := newDeploy.GetProject()
		fmt.Printf("Et bien %v\n", project)
		fmt.Sprint(project)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "Deployed PR will comment all PR deployed to %q", newDeploy.Target)
		project.GetClosedPullRequests()
	} else {
		http.Error(w, "You must send your request in POST.", 405)
	}
}

func main() {
	http.HandleFunc("/", processRequest)
	http.ListenAndServe(":8080", nil)
}
