package philaeprobe

import (
	"encoding/json"
	"net/http"

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

func (p PhilaeProbe) Check() bool {
	client := &http.Client{}

	req, err := http.NewRequest("GET", p.endpoint, nil)
	if err != nil {
		return false
	}

	if p.user != "" || p.password != "" {
		req.SetBasicAuth(p.user, p.password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}

	if resp.Status[0] != '2' && resp.Status[0] != '3' {
		return false
	}

	var result prober.Result

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false
	}

	return result.Healthy
}
