package deptools

import (
	"bytes"
	"net/http"
	"testing"
)

func TestRequestAnalyserWithoutOwnerInRequestReturnFalse(t *testing.T) {
	postJson := []byte(`{"Repo": "carenelec-y3b", "BaseType": "branch", "BaseName": "master", "Target": "Prod"}`)
	request, _ := http.NewRequest("POST", "/", bytes.NewReader(postJson))
	var requestAnalyser RequestAnalyser
	_, _, _, _, _, parseError := requestAnalyser.Parse(request)
	if parseError == nil {
		t.Error("isValid() on request without Owner should return false")
	}
}

func TestRequestAnalyserWithoutRepoInRequestReturnFalse(t *testing.T) {
	postJson := []byte(`{"Owner":"alexisjanvier", "BaseType": "branch", "BaseName": "master", "Target": "Prod"}`)
	request, _ := http.NewRequest("POST", "/", bytes.NewReader(postJson))
	var requestAnalyser RequestAnalyser
	_, _, _, _, _, parseError := requestAnalyser.Parse(request)
	if parseError == nil {
		t.Error("isValid() on request without Repo should return false")
	}
}

func TestRequestAnalyserWithoutBaseTypeInRequestReturnFalse(t *testing.T) {
	postJson := []byte(`{"Owner":"alexisjanvier", "Repo": "carenelec-y3b", "BaseName": "master", "Target": "Prod"}`)
	request, _ := http.NewRequest("POST", "/", bytes.NewReader(postJson))
	var requestAnalyser RequestAnalyser
	_, _, _, _, _, parseError := requestAnalyser.Parse(request)
	if parseError == nil {
		t.Error("isValid() on request without BaseType should return false")
	}
}

func TestRequestAnalyserWithoutBaseNameInRequestReturnFalse(t *testing.T) {
	postJson := []byte(`{"Owner":"alexisjanvier", "Repo": "carenelec-y3b", "BaseType": "branch", "Target": "Prod"}`)
	request, _ := http.NewRequest("POST", "/", bytes.NewReader(postJson))
	var requestAnalyser RequestAnalyser
	_, _, _, _, _, parseError := requestAnalyser.Parse(request)
	if parseError == nil {
		t.Error("isValid() on request without BaseName should return false")
	}
}

func TestRequestAnalyserWithoutTargetInRequestReturnFalse(t *testing.T) {
	postJson := []byte(`{"Owner":"alexisjanvier", "Repo": "carenelec-y3b", "BaseType": "branch", "BaseName": "master"}`)
	request, _ := http.NewRequest("POST", "/", bytes.NewReader(postJson))
	var requestAnalyser RequestAnalyser
	_, _, _, _, _, parseError := requestAnalyser.Parse(request)
	if parseError == nil {
		t.Error("isValid() on request without Target should return false")
	}
}
