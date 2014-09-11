package deployment

type PullRequest struct {
	Number  int
	Title   string
	HeadRef string
	HeadSHA string
	BaseRef string
}

func (pr *PullRequest) isMergedOnBranch(branch string) bool {
	return true
}

func (pr *PullRequest) isMergedOnTag(tag string) bool {
	return true
}

func (pr *PullRequest) isAlreadyDeployOnTarget(target string) bool {
	return true
}

func (pr *PullRequest) commentAsDeployOn(target string) bool {
	return true
}
