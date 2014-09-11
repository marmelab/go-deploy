package deployment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PullRequest struct {
	Number   int
	Title    string
	HeadRef  string
	HeadSHA  string
	Status   string
	MergeSHA string
}

type jsonPullRequest struct {
	MergeCommitSHA string `json:"merge_commit_sha"`
}

func (pr *PullRequest) IsMergedOnBranch(branch string) bool {
	return true
}

func (pr *PullRequest) IsMergedOnTag(tag string) bool {
	return true
}

func (pr *PullRequest) IsAlreadyDeployOnTarget(target string) bool {
	return true
}

func (pr *PullRequest) CommentAsDeployOn(target string) bool {
	return true
}

func (pr *PullRequest) getMergeSHA(owner string, repo string) {

	gitHubCall := fmt.Sprintf("https://api.github.com/repos/%v/%v/pulls/%v", owner, repo, pr.Number)
	//TODO error managing
	res, err := http.Get(gitHubCall)
	if err != nil {
		fmt.Printf("error: %v\n\n", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("error: %v\n\n", err)
	}

	if res.StatusCode != 200 {
		fmt.Printf("error: %v\n\n", res.StatusCode)
	}

	jsonPr := jsonPullRequest{}

	err = json.Unmarshal(body, &jsonPr)
	if err != nil {
		fmt.Printf("error: %v\n\n", err)
	}

	pr.MergeSHA = jsonPr.MergeCommitSHA
}
