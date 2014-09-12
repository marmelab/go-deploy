package deptools

type Project struct {
	Owner       string
	Repo        string
	AccessToken string
}

func (project *Project) IsConfig() error {
	// TODO make a config.json with all configured projects with their API Keys
	project.AccessToken = "4a70688e761dc2280eb9dec5bb833807709c69eb"
	return nil
}
