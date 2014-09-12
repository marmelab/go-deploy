package deptools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Project struct {
	Owner       string
	Repo        string
	AccessToken string
}

type Configuration struct {
	Projects []ProjectConf
}

type ProjectConf struct {
	Owner       string
	Repository  string
	AccessToken string
}

func (project *Project) IsConfig() error {
	file, readFileError := ioutil.ReadFile("./config.json")
	if readFileError != nil {
		return readFileError
	}
	var config Configuration
	readJsonError := json.Unmarshal(file, &config)
	if readJsonError != nil {
		return readJsonError
	}
	for _, projectConf := range config.Projects {
		if projectConf.Owner == project.Owner && projectConf.Repository == project.Repo {
			project.AccessToken = projectConf.AccessToken
			return nil
		}
	}
	err := fmt.Errorf("project %v/%v is not present in the config.json file", project.Owner, project.Repo)
	return err
}
