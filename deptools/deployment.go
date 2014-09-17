package deptools

import (
	"fmt"
	"github.com/google/go-github/github"
	"time"
)

const (
	BASE_TAG    = "tag"
	BASE_BRANCH = "branch"
)

type Deployment struct {
	Owner                       string                 `bson:"Owner"`
	Repository                  string                 `bson:"Repository"`
	Base_type                   string                 `bson:"base_type"`
	Base_name                   string                 `bson:"base_name"`
	Target                      string                 `bson:"Target"`
	Created_at                  time.Time              `bson:"created_at"`
	Github_client               *github.Client         `bson:"-"`
	Last_pr_merge_date          time.Time              `bson:"last_pr_merge_date"`
	base_tag_SHA                string                 `bson:"-"`
	pull_requests               map[string]PullRequest `bson:"-"`
	commits_on_deployed_base    map[string]string      `bson:"-"`
	pull_request_merged_on_base map[string]PullRequest `bson:"-"`
}

type PullRequest struct {
	Number   int       `bson:"number"`
	Title    string    `bson:"-"`
	HeadRef  string    `bson:"-"`
	HeadSHA  string    `bson:"header_sha"`
	Status   string    `bson:"-"`
	MergedAt time.Time `bson:"merged_at"`
}

type PrCommentedToTarget struct {
	Number  int    `bson:"number"`
	HeadSHA string `bson:"header_sha"`
	Target  string `bson:"target"`
}

func (dpl *Deployment) BaseExist() (bool, error) {
	if dpl.Base_type == BASE_TAG {
		exist, err := dpl.tagExist()
		return exist, err
	} else {
		exist, err := dpl.branchExist()
		return exist, err
	}
}

//TODO refactoring these two identical functions by regulating the type signature issues
func (dpl *Deployment) tagExist() (bool, error) {
	opt := &github.ListOptions{}
	tags, _, err := dpl.Github_client.Repositories.ListTags(dpl.Owner, dpl.Repository, opt)
	if err != nil {
		return false, err
	}
	for _, tag := range tags {
		if *tag.Name == dpl.Base_name {
			dpl.base_tag_SHA = *tag.Commit.SHA
			return true, nil
		}
	}

	tagNotFound := fmt.Errorf("%v %v not found on %v", dpl.Base_type, dpl.Base_name, dpl.Repository)
	return false, tagNotFound
}
func (dpl *Deployment) branchExist() (bool, error) {
	opt := &github.ListOptions{}
	branches, _, err := dpl.Github_client.Repositories.ListBranches(dpl.Owner, dpl.Repository, opt)
	if err != nil {
		return false, err
	}
	for _, branch := range branches {
		if *branch.Name == dpl.Base_name {
			return true, nil
		}
	}

	branchNotFound := fmt.Errorf("%v %v not found on %v", dpl.Base_type, dpl.Base_name, dpl.Repository)
	return false, branchNotFound
}

func (dpl *Deployment) CommentPrContainedInDeploy() (int, error) {
	if fetchPRError := dpl.fetchClosedPullRequests(); fetchPRError != nil {
		return 0, fetchPRError
	}

	if fetchCommitError := dpl.fetchCommitsOnBase(); fetchCommitError != nil {
		return 0, fetchCommitError
	}

	dpl.setPrMergedOnBase()
	if nbPrToComment := len(dpl.pull_request_merged_on_base); nbPrToComment < 1 {
		return 0, nil
	}
	nbPrCommented, commentError := dpl.commentMergedPR()
	if commentError != nil {
		return 0, commentError
	}

	return nbPrCommented, nil
}

func (dpl *Deployment) fetchClosedPullRequests() error {

	//TODO limit the number of result, with Base Option ?
	//TODO refactoring because if else -> for -> if ....
	opt := &github.PullRequestListOptions{State: "closed"}
	prs, _, getPrError := dpl.Github_client.PullRequests.List(dpl.Owner, dpl.Repository, opt)
	if getPrError != nil {
		return getPrError
	}
	dpl.pull_requests = make(map[string]PullRequest)
	for _, pr := range prs {
		if pr.MergedAt == nil {
			continue
		}
		pullr := PullRequest{
			Number:   *pr.Number,
			Title:    *pr.Title,
			HeadRef:  *pr.Head.Ref,
			HeadSHA:  *pr.Head.SHA,
			Status:   *pr.State,
			MergedAt: *pr.MergedAt}
		dpl.pull_requests[pullr.HeadSHA] = pullr
	}

	return nil
}

func (dpl *Deployment) fetchCommitsOnBase() error {
	lastPrMergedDate := dpl.getLastPrMergeDate()
	var searchBaseName string
	if dpl.Base_type == "tag" {
		searchBaseName = string(dpl.base_tag_SHA)
	} else {
		searchBaseName = string(dpl.Base_name)
	}
	opt := &github.CommitsListOptions{SHA: searchBaseName, Since: lastPrMergedDate}
	commits, _, err := dpl.Github_client.Repositories.ListCommits(dpl.Owner, dpl.Repository, opt)
	if err != nil {
		return err
	}
	dpl.commits_on_deployed_base = make(map[string]string)
	for _, commit := range commits {
		dpl.commits_on_deployed_base[*commit.SHA] = *commit.SHA
	}

	return nil
}

func (dpl *Deployment) setPrMergedOnBase() {
	dpl.pull_request_merged_on_base = make(map[string]PullRequest)
	for commitSha, _ := range dpl.commits_on_deployed_base {
		if mergedPullRequest, found := dpl.pull_requests[commitSha]; found {
			dpl.pull_request_merged_on_base[mergedPullRequest.HeadSHA] = mergedPullRequest
			if mergedPullRequest.MergedAt.After(dpl.Last_pr_merge_date) {
				dpl.Last_pr_merge_date = mergedPullRequest.MergedAt
			}
			fmt.Println(".")
		}
	}
}

func (dpl *Deployment) commentMergedPR() (int, error) {
	msg := fmt.Sprintf("This PR was deployed to %v (from the %v %v)", dpl.Target, dpl.Base_type, dpl.Base_name)
	comment := &github.IssueComment{Body: &msg}
	nbPrCommented := 0
	for _, prToComment := range dpl.pull_request_merged_on_base {
		if hasAlreadyBeenCommented := prToComment.hasBeenDeployTo(dpl.Target); !hasAlreadyBeenCommented {
			_, _, err := dpl.Github_client.Issues.CreateComment(dpl.Owner, dpl.Repository, prToComment.Number, comment)
			if err != nil {
				return nbPrCommented, err
			}
			prToComment.saveAsCommentToTarget(dpl.Target)
			nbPrCommented += 1
		}
	}
	dpl.save()

	return nbPrCommented, nil
}
