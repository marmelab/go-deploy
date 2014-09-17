package main

import (
	"bitbucket.org/alexisjanvier/deployedpr/deptools"
	"fmt"
	"net/http"
	"time"
)

func processRequest(response http.ResponseWriter, r *http.Request) {
	defer func() {
		if panicError := recover(); panicError != nil {
			errorMsg := fmt.Sprintf("%q", panicError)
			http.Error(response, errorMsg, 400)
		}
	}()
	//TODO accept also GET request
	if r.Method != "POST" {
		http.Error(response, "You must send your request in POST.", 405)
		return
	}

	var requestAnalyser deptools.RequestAnalyser
	owner, repo, baseType, baseName, target, parseError := requestAnalyser.Parse(r)
	if parseError != nil {
		panic(parseError)
	}

	gitHubClient, ghClientError := deptools.GetGithubClient(owner, repo)
	if ghClientError != nil {
		panic(ghClientError)
	}

	deploy := deptools.Deployment{
		Owner:         owner,
		Repository:    repo,
		Base_type:     baseType,
		Base_name:     baseName,
		Target:        target,
		Created_at:    time.Now(),
		Github_client: gitHubClient,
	}
	if baseExist, baseError := deploy.BaseExist(); !baseExist || baseError != nil {
		panic(baseError)
	}
	nbPrCommented, commentError := deploy.CommentPrContainedInDeploy()
	if commentError != nil {
		panic(nil)
	}
	if nbPrCommented < 1 {
		fmt.Fprintf(response, "There was no PR to comment on this deployment to %q", target)
	} else {
		fmt.Fprintf(response, "%v PR contained in this deployment on %q have been commented", nbPrCommented, target)
	}

	return
}

func main() {
	http.HandleFunc("/", processRequest)
	http.ListenAndServe(":8080", nil)
}
