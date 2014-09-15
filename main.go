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
	if r.Method != "POST" {
		http.Error(response, "You must send your request in POST.", 405)
		return
	}
	var requestAnalyser deptools.RequestAnalyser
	owner, repo, baseType, baseName, target, parseError := requestAnalyser.Parse(r)
	if parseError != nil {
		panic(parseError)
	}

	project := deptools.Project{Owner: owner, Repo: repo}
	if projectConfigError := project.IsConfig(); projectConfigError != nil {
		panic(projectConfigError)
	}

	deploy := deptools.Deployment{
		Owner:       owner,
		Repository:  repo,
		AccessToken: project.AccessToken,
		BaseType:    baseType,
		BaseName:    baseName,
		Target:      target,
		CreatedAt:   time.Now(),
	}
	if baseExist, baseError := deploy.BaseExist(); !baseExist || baseError != nil {
		panic(baseError)
	}
	pullRequestsCommented, commentError := deploy.CommentPrContainedInDeploy()
	if commentError != nil {
		panic(nil)
	}
	if pullRequestsCommented == "" {
		fmt.Fprintf(response, "There was no PR to comment on this deployment to %q", target)
	} else {
		fmt.Fprintf(response, "The PR contained in this deployment on %q have been commented", target)
	}

	return
}

func main() {
	http.HandleFunc("/", processRequest)
	http.ListenAndServe(":8080", nil)
}
