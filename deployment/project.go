package deployment

import (
	"github.com/google/go-github/github"
)

type Project struct {
	Owner string
	Repo  string
}

func (project *Project) getClosedPullRequests() bool {
	client := github.NewClient(nil)
	opt := &github.PullRequestListOptions{State: "closed"}
	prs, _, err := client.PullRequests.List(project.Owner, project.Repo, opt)
	if err != nil {
		return true
	} else {
		return true
	}
}
