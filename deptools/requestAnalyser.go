package deptools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestAnalyser struct {
	Owner    string
	Repo     string
	BaseType string
	BaseName string
	Target   string
}

func (ra *RequestAnalyser) Parse(r *http.Request) (owner string, repo string, basetype string, basename string, target string, analyseError error) {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&ra); err != nil {
		return "", "", "", "", "", errors.New("Request as not a valid json format")
	}
	if jsonValid, erroMsg := ra.IsValid(); !jsonValid {
		return "", "", "", "", "", errors.New(erroMsg)
	}
	return ra.Owner,
		ra.Repo,
		ra.BaseType,
		ra.BaseName,
		ra.Target,
		nil

}

func (ra *RequestAnalyser) IsValid() (bool, string) {
	isValid := true
	errorMsg := ""
	if ra.Owner == "" {
		isValid = false
		errorMsg = fmt.Sprintf("Owner must be set in your json %s", errorMsg)
	}
	if ra.Repo == "" {
		isValid = false
		errorMsg = fmt.Sprintf("Repo must be set in your json %s", errorMsg)
	}
	if ra.BaseType == "" {
		isValid = false
		errorMsg = fmt.Sprintf("BaseType must be set in your json %s", errorMsg)
	}
	if ra.BaseName == "" {
		isValid = false
		errorMsg = fmt.Sprintf("BaseName must be set in your json, not both %s", errorMsg)
	}
	if ra.Target == "" {
		isValid = false
		errorMsg = fmt.Sprintf("Target must be set in your json, not both %s", errorMsg)
	}

	return isValid, errorMsg
}
