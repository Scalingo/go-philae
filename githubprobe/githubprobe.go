package githubprobe

import (
	"context"
	"encoding/json"
	"io"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/go-philae/v5/httpprobe"
)

type GithubChecker struct{}

type GithubProbe struct {
	http httpprobe.HTTPProbe
}

func (_ GithubChecker) Check(body io.Reader) error {
	var result GithubStatusResponse

	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		return errgo.Notef(err, "invalid json")
	}

	if result.Status.Indicator != "none" {
		return errgo.Newf("GitHub is probably down")
	}
	return nil
}

func NewGithubProbe(name string) GithubProbe {
	return GithubProbe{
		http: httpprobe.NewHTTPProbe(name, "https://www.githubstatus.com/api/v2/status.json", httpprobe.HTTPOptions{
			Checker: GithubChecker{},
		}),
	}
}

func (p GithubProbe) Name() string {
	return p.http.Name()
}

func (p GithubProbe) Check(ctx context.Context) error {
	return p.http.Check(ctx)
}
