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
	BaseTagSHA              string
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
	if dpl.BaseType == "tag" {
		exist, err := dpl.tagExist()
		return exist, err
	} else {
		exist, err := dpl.branchExist()
		return exist, err
	}
}

//TODO refactoring these two identical functions by regulating the type signature issues
func (dpl *Deployment) tagExist() (bool, error) {
	client := dpl.getGithubAccessClient()
	opt := &github.ListOptions{}
	tags, _, err := client.Repositories.ListTags(dpl.Project.Owner, dpl.Project.Repo, opt)
	if err != nil {
		return false, err
	} else {
		for _, tag := range tags {
			if *tag.Name == dpl.BaseName {
				dpl.BaseTagSHA = *tag.Commit.SHA
				return true, nil
			}
		}
	}
	tagNotFound := fmt.Errorf("%v %v not found on %v", dpl.BaseType, dpl.BaseName, dpl.Project.Repo)
	return false, tagNotFound
}
func (dpl *Deployment) branchExist() (bool, error) {
	client := dpl.getGithubAccessClient()
	opt := &github.ListOptions{}
	branches, _, err := client.Repositories.ListBranches(dpl.Project.Owner, dpl.Project.Repo, opt)
	if err != nil {
		return false, err
	} else {
		for _, branch := range branches {
			if *branch.Name == dpl.BaseName {
				return true, nil
			}
		}
	}
	branchNotFound := fmt.Errorf("%v %v not found on %v", dpl.BaseType, dpl.BaseName, dpl.Project.Repo)
	return false, branchNotFound
}

func (dpl *Deployment) CommentPrContainedInDeploy() (string, error) {
	if getPRError := dpl.getClosedPullRequests(); getPRError != nil {
		return "", getPRError
	}

	if getCommitError := dpl.getCommitsOnBase(); getCommitError != nil {
		return "", getCommitError
	}

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
	client := dpl.getGithubAccessClient()
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
	var searchBaseName string
	if dpl.BaseType == "tag" {
		searchBaseName = string(dpl.BaseTagSHA)
	} else {
		searchBaseName = string(dpl.BaseName)
	}
	client := dpl.getGithubAccessClient()
	opt := &github.CommitsListOptions{SHA: searchBaseName}
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
	client := dpl.getGithubAccessClient()
	msg := fmt.Sprintf("This pull Request as been deployed on %v (from the %v %v)", dpl.Target, dpl.BaseType, dpl.BaseName)
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

func (dpl *Deployment) getGithubAccessClient() *github.Client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: dpl.Project.AccessToken},
	}
	client := github.NewClient(t.Client())
	return client
}
