package main

import (
	"bitbucket.org/alexisjanvier/deployedpr/deptools"
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
		var requestAnalyser deptools.RequestAnalyser
		owner, repo, baseType, baseName, target, parseError := requestAnalyser.Parse(r)
		if parseError != nil {
			panic(parseError)
		}

		project := deptools.Project{Owner: owner, Repo: repo}
		if projectConfigError := project.IsConfig(); projectConfigError != nil {
			panic(projectConfigError)
		}

		deploy := deptools.Deployment{Project: project, BaseType: baseType, BaseName: baseName, Target: target}
		if baseExist, baseError := deploy.BaseExist(); !baseExist || baseError != nil {
			panic(baseError)
		}
		pullRequestsCommented, commentError := deploy.CommentPrContainedInDeploy()
		if commentError != nil {
			panic(nil)
		}
		if pullRequestsCommented == "" {
			fmt.Fprintf(w, "There were no PR to comment on this deployment to %q", target)
		}

		fmt.Fprintf(w, "Deployed PR will comment all PR deployed to %q", target)
	} else {
		http.Error(w, "You must send your request in POST.", 405)
	}
}

func main() {
	http.HandleFunc("/", processRequest)
	http.ListenAndServe(":8080", nil)
}
