package githubprobe

type GithubStatusResponse struct {
	Status GithubStatusResponseStatus `json:"status"`
}

type GithubStatusResponseStatus struct {
	Indicator string `json:"indicator"`
}
