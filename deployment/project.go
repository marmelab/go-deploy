package deployment

import (
	"fmt"
	"github.com/google/go-github/github"
)

type Project struct {
	Owner string
	Repo  string
}

func (project *Project) GetClosedPullRequests() bool {
	client := github.NewClient(nil)
	opt := &github.PullRequestListOptions{State: "closed"}
	prs, _, err := client.PullRequests.List(project.Owner, project.Repo, opt)
	if err != nil {
		fmt.Printf("error: %v\n\n", err)
	} else {
		for _, pr := range prs {
			fmt.Println("\n\n\n############")
			fmt.Printf("Number: %v\n", github.Stringify(pr.Number))
			fmt.Printf("State: %v\n", github.Stringify(pr.State))
			fmt.Printf("Title: %v\n", github.Stringify(pr.Title))
			fmt.Printf("Head label: %v\n", github.Stringify(pr.Head.Label))
			fmt.Printf("Head sha: %v\n", github.Stringify(pr.Head.SHA))
			fmt.Printf("Base label: %v", github.Stringify(pr.Base.Label))
		}
	}

	return true
}
