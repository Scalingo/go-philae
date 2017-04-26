package philaeprobe

import (
	"encoding/json"
	"io"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/go-philae/httpprobe"
	"github.com/Scalingo/go-philae/prober"
)

type PhilaeProbe struct {
	http httpprobe.HTTPProbe
}

type PhilaeChecker struct{}

func (_ PhilaeChecker) Check(body io.Reader) error {
	var result prober.Result

	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		return errgo.Notef(err, "Invalid json")
	}

	if !result.Healthy {
		return errgo.Newf("Node not healthy")
	}

	return nil
}

func NewPhilaeProbe(name, endpoint string) PhilaeProbe {
	return PhilaeProbe{
		http: httpprobe.NewCheckedHTTPProbe(name, endpoint, PhilaeChecker{}),
	}
}

func NewAuthenticatedPhilaeProbe(name, endpoint, user, password string) PhilaeProbe {
	return PhilaeProbe{
		http: httpprobe.NewAuthenticatedCheckedHTTPProbe(name, endpoint, user, password, PhilaeChecker{}),
	}
}

func (p PhilaeProbe) Name() string {
	return p.http.Name()
}

func (p PhilaeProbe) Check() error {
	return p.http.Check()
}
