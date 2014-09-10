package deployment

import (
	"testing"
)

func TestDeploymentWithoutGitRefIsNotValid(t *testing.T) {
	deploy := Deployment{Branch: "master", Target: "prod"}
	if jsonValid, _ := deploy.IsValid(); jsonValid {
		t.Error("isValid() on Deployment without GitRef should return false")
	}
}

func TestDeploymentWithoutTargetIsNotValid(t *testing.T) {
	deploy := Deployment{GitRef: "alexisjanvier/dummy-project", Branch: "master"}
	if jsonValid, _ := deploy.IsValid(); jsonValid {
		t.Error("isValid() on Deployment without Target should return false")
	}
}

func TestDeploymentWithoutBranchOrTagIsNotValid(t *testing.T) {
	deploy := Deployment{GitRef: "alexisjanvier/dummy-project", Target: "prod"}
	if jsonValid, _ := deploy.IsValid(); jsonValid {
		t.Error("isValid() on Deployment without Branch or Tag should return false")
	}
}

func TestDeploymentWithBranchAndTagIsNotValid(t *testing.T) {
	deploy := Deployment{GitRef: "alexisjanvier/dummy-project", Branch: "master", Tag: "v1", Target: "prod"}
	if jsonValid, _ := deploy.IsValid(); jsonValid {
		t.Error("isValid() on Deployment with Branch and Tag should return false")
	}
}
