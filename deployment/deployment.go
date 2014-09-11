package deployment

import (
	"errors"
	"fmt"
	"strings"
)

type Deployment struct {
	GitRef string
	Branch string
	Tag    string
	Target string
}

func (dpl *Deployment) IsValid() (isValid bool, errorMsg string) {
	isValid = true
	errorMsg = ""
	if dpl.GitRef == "" {
		isValid = false
		errorMsg = fmt.Sprintf("GitRef must be set in your json %s", errorMsg)
	}
	if dpl.Target == "" {
		isValid = false
		errorMsg = fmt.Sprintf("Target must be set in your json %s", errorMsg)
	}
	if dpl.Branch == "" && dpl.Tag == "" {
		isValid = false
		errorMsg = fmt.Sprintf("A branch or a target must be set in your json %s", errorMsg)
	}
	if dpl.Branch != "" && dpl.Tag != "" {
		isValid = false
		errorMsg = fmt.Sprintf("A branch or a target must be set in your json, not both %s", errorMsg)
	}
	return
}

func (dpl *Deployment) GetProject() (project Project, errorMsg error) {
	owner, repo, err := dpl.splitGitRef()
	if err != nil {
		return Project{Owner: "", Repo: ""}, err
	}
	project = Project{Owner: owner, Repo: repo}

	return project, nil
}

func (dpl *Deployment) splitGitRef() (owner string, repo string, err error) {
	s := strings.Split(dpl.GitRef, "/")
	if len(s) < 2 {
		return "", "", errors.New("GitRef must be format as owner/repo")
	}
	return s[0], s[1], nil
}
