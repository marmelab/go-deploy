package deployment

import (
	"fmt"
	"github.com/google/go-github/github"
)

type Project struct {
	Owner        string
	Repo         string
	PullRequests map[int]PullRequest
}

func (project *Project) GetClosedPullRequests() {
	client := github.NewClient(nil)
	opt := &github.PullRequestListOptions{State: "closed"}
	prs, _, err := client.PullRequests.List(project.Owner, project.Repo, opt)
	if err != nil {
		//TODO use panic()
		fmt.Printf("error: %v\n\n", err)
	} else {
		//fmt.Printf("State: %v\n", *prs[0].Number)
		project.PullRequests = make(map[int]PullRequest)
		for _, pr := range prs {
			pullr := PullRequest{
				Number:  *pr.Number,
				Title:   *pr.Title,
				HeadRef: *pr.Head.Ref,
				HeadSHA: *pr.Head.SHA,
				Status:  *pr.State}
			pullr.getMergeSHA(project.Owner, project.Repo)
			project.PullRequests[pullr.Number] = pullr
		}
	}

	return
}
