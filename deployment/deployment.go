package deployment

import (
	"fmt"
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
