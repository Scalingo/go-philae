package gitlabprobe

import (
	"encoding/json"
	"io"

	"gopkg.in/errgo.v1"

	"github.com/Scalingo/go-philae/httpprobe"
	"github.com/Scalingo/go-philae/statusioprobe"
)

type GitlabProbe struct {
	http httpprobe.HTTPProbe
}

type GitlabChecker struct{}

func NewGitlabProbe(name string) GitlabProbe {
	return GitlabProbe{
		// ID get from here : https://status.gitlab.com in the HTTP Header : x-status-page-id
		http: httpprobe.NewHTTPProbe(name, "https://api.status.io/1.0/status/5b36dc6502d06804c08349f7", httpprobe.HTTPOptions{
			Checker: GitlabChecker{},
		}),
	}
}

func (p GitlabProbe) Name() string {
	return p.http.Name()
}

func (p GitlabProbe) Check() error {
	return p.http.Check()
}

func (GitlabChecker) Check(body io.Reader) error {
	var result statusioprobe.StatusIOResponse

	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		return errgo.Notef(err, "Invalid json")
	}

	if result.Result.Overall.StatusCode >= 400 {
		return errgo.Newf("One or more services from GitLab are down!")
	}

	return nil
}
