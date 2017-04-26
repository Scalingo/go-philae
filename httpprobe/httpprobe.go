package httpprobe

import (
	"net/http"

	errgo "gopkg.in/errgo.v1"
)

type HTTPProbe struct {
	name     string
	endpoint string
	user     string
	password string
}

func NewHTTPProbe(name, endpoint string) HTTPProbe {
	return HTTPProbe{
		name:     name,
		endpoint: endpoint,
		user:     "",
		password: "",
	}
}

func NewAuthenticatedHTTPProbe(name, endpoint, user, password string) HTTPProbe {
	return HTTPProbe{
		name:     name,
		endpoint: endpoint,
		user:     user,
		password: password,
	}
}

func (p HTTPProbe) Name() string {
	return p.name
}

func (p HTTPProbe) Check() error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", p.endpoint, nil)
	if err != nil {
		return errgo.Notef(err, "Unable to create request")
	}

	if p.user != "" || p.password != "" {
		req.SetBasicAuth(p.user, p.password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return errgo.Notef(err, "Unable to send request")
	}

	if resp.Status[0] != '2' && resp.Status[0] != '3' {
		return errgo.Newf("Invalid return code: %s", resp.Status)
	}

	return nil
}
