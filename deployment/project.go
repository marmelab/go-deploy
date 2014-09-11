package deployment

import (
	"fmt"
	"github.com/google/go-github/github"
)

type Project struct {
	Owner                 string
	Repo                  string
	PullRequests          map[string]PullRequest
	CommitsOnDeployBranch map[string]string
}

func (project *Project) GetClosedPullRequests() {
	client := github.NewClient(nil)
	opt := &github.PullRequestListOptions{State: "closed"}
	prs, _, err := client.PullRequests.List(project.Owner, project.Repo, opt)
	if err != nil {
		//TODO use panic()
		fmt.Printf("error: %v\n\n", err)
	} else {
		project.PullRequests = make(map[string]PullRequest)
		for _, pr := range prs {
			pullr := PullRequest{
				Number:  *pr.Number,
				Title:   *pr.Title,
				HeadRef: *pr.Head.Ref,
				HeadSHA: *pr.Head.SHA,
				Status:  *pr.State}
			//pullr.getMergeSHA(project.Owner, project.Repo)
			project.PullRequests[pullr.HeadSHA] = pullr
		}
	}

	return
}

func (project *Project) GetCommitsOnBranch(branch string) {
	client := github.NewClient(nil)
	opt := &github.CommitsListOptions{SHA: branch}
	commits, _, err := client.Repositories.ListCommits(project.Owner, project.Repo, opt)
	if err != nil {
		//TODO use panic()
		fmt.Printf("error: %v\n\n", err)
	} else {
		project.CommitsOnDeployBranch = make(map[string]string)
		for _, commit := range commits {
			project.CommitsOnDeployBranch[*commit.SHA] = *commit.SHA
		}
	}

}

func (project *Project) GetPrMergedOnBranch() {
	for commitSha, _ := range project.CommitsOnDeployBranch {
		if mergedPullRequest, found := project.PullRequests[commitSha]; found {
			fmt.Printf("PR %v (%v) has been merged in branch\n", mergedPullRequest.Number, mergedPullRequest.Title)
		}
	}
}
