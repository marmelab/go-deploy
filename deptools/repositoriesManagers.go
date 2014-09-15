package deptools

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"io/ioutil"
)

type Configuration struct {
	Projects []ProjectConf `json:"projects"`
}

type ProjectConf struct {
	Owner       string `json:"owner"`
	Repository  string `json:"repository"`
	AccessToken string `json:"access_token"`
}

func GetGithubClient(owner string, repository string) (*github.Client, error) {
	accessToken, getTokenError := loadAccessToken(owner, repository)
	if getTokenError != nil {

		return github.NewClient(nil), getTokenError
	}
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: accessToken},
	}

	return github.NewClient(t.Client()), nil
}

func loadAccessToken(owner string, repository string) (string, error) {
	file, readFileError := ioutil.ReadFile("./config.json")
	if readFileError != nil {

		return "", readFileError
	}
	var config Configuration
	readJsonError := json.Unmarshal(file, &config)
	if readJsonError != nil {

		return "", readJsonError
	}
	for _, project := range config.Projects {
		if project.Owner == owner && project.Repository == repository {

			return project.AccessToken, nil
		}
	}
	projectConfigError := fmt.Errorf("project %v/%v is not present in the config.json file", owner, repository)

	return "", projectConfigError
}
