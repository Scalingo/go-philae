package philaeprobe

import (
	"encoding/json"
	"net/http"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/go-philae/prober"
)

type PhilaeProbe struct {
	name     string
	endpoint string
	user     string
	password string
}

func NewPhilaeProbe(name, endpoint string) PhilaeProbe {
	return PhilaeProbe{
		name:     name,
		endpoint: endpoint,
		user:     "",
		password: "",
	}
}

func NewAuthenticatedPhilaeProbe(name, endpoint, user, password string) PhilaeProbe {
	return PhilaeProbe{
		name:     name,
		endpoint: endpoint,
		user:     user,
		password: password,
	}
}

func (p PhilaeProbe) Name() string {
	return p.name
}

func (p PhilaeProbe) Check() error {
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
		return errgo.Notef(err, "Invalid return code")
	}

	var result prober.Result

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return errgo.Notef(err, "Invalid json")
	}

	if !result.Healthy {
		return errgo.Newf("Node not healthy")
	}

	return nil
}
