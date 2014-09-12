package deptools

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/google/go-github/github"
)

type Deployment struct {
	Project                 Project
	BaseType                string
	BaseName                string
	Target                  string
	PullRequests            map[string]PullRequest
	CommitsOnDeployedBase   map[string]string
	PullRequestMergedOnBase map[string]int
}

type PullRequest struct {
	Number  int
	Title   string
	HeadRef string
	HeadSHA string
	Status  string
}

func (dpl *Deployment) BaseExist() (bool, error) {
	// TODO Check if tag or branch exist in repo
	return true, nil
}

func (dpl *Deployment) CommentPrContainedInDeploy() (string, error) {
	if getPRError := dpl.getClosedPullRequests(); getPRError != nil {
		return "", getPRError
	}

	if getCommitError := dpl.getCommitsOnBase(); getCommitError != nil {
		return "", getCommitError
	}
	// 3)On croise
	dpl.setPrMergedOnBase()
	if nbPrToComment := len(dpl.PullRequestMergedOnBase); nbPrToComment < 1 {
		return "", nil
	}
	if commentError := dpl.commentMergedPR(); commentError != nil {
		return "", commentError
	}
	return "DONE", nil
}

func (dpl *Deployment) getClosedPullRequests() error {
	client := github.NewClient(nil)
	opt := &github.PullRequestListOptions{State: "closed"}
	prs, _, err := client.PullRequests.List(dpl.Project.Owner, dpl.Project.Repo, opt)
	if err != nil {
		return err
	} else {
		dpl.PullRequests = make(map[string]PullRequest)
		for _, pr := range prs {
			pullr := PullRequest{
				Number:  *pr.Number,
				Title:   *pr.Title,
				HeadRef: *pr.Head.Ref,
				HeadSHA: *pr.Head.SHA,
				Status:  *pr.State}
			dpl.PullRequests[pullr.HeadSHA] = pullr
		}
	}
	return nil
}

func (dpl *Deployment) getCommitsOnBase() error {
	//TODO make case where baseType is tag
	client := github.NewClient(nil)
	opt := &github.CommitsListOptions{SHA: dpl.BaseName}
	commits, _, err := client.Repositories.ListCommits(dpl.Project.Owner, dpl.Project.Repo, opt)
	if err != nil {
		return err
	} else {
		dpl.CommitsOnDeployedBase = make(map[string]string)
		for _, commit := range commits {
			dpl.CommitsOnDeployedBase[*commit.SHA] = *commit.SHA
		}
	}
	return nil
}

func (dpl *Deployment) setPrMergedOnBase() {
	dpl.PullRequestMergedOnBase = make(map[string]int)
	for commitSha, _ := range dpl.CommitsOnDeployedBase {
		if mergedPullRequest, found := dpl.PullRequests[commitSha]; found {
			dpl.PullRequestMergedOnBase[mergedPullRequest.HeadSHA] = mergedPullRequest.Number
		}
	}
}

func (dpl *Deployment) commentMergedPR() error {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: dpl.Project.AccessToken},
	}
	client := github.NewClient(t.Client())
	msg := fmt.Sprintf("This pull Request as been deployed on %v", dpl.Target)
	comment := &github.IssueComment{Body: &msg}
	for _, prNumber := range dpl.PullRequestMergedOnBase {
		// TODO get created comment and save it somewhere commentPr, _, err
		_, _, err := client.Issues.CreateComment(dpl.Project.Owner, dpl.Project.Repo, prNumber, comment)
		if err != nil {
			return err
		}
	}
	return nil
}
