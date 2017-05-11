package httpprobe

import (
	"net"
	"net/http"
	"time"

	errgo "gopkg.in/errgo.v1"
)

type HTTPProbe struct {
	name     string
	endpoint string
	options  HTTPOptions
}

type HTTPOptions struct {
	Username           string
	Password           string
	Checker            HTTPChecker
	ExpectedStatusCode int
	testing            bool
}

func NewHTTPProbe(name, endpoint string, opts HTTPOptions) HTTPProbe {
	return HTTPProbe{
		name:     name,
		endpoint: endpoint,
		options:  opts,
	}
}

func (p HTTPProbe) Name() string {
	return p.name
}

func (p HTTPProbe) Check() error {
	client := NewTimeoutClient()

	if p.options.testing {
		client = &http.Client{}
	}

	req, err := http.NewRequest("GET", p.endpoint, nil)
	if err != nil {
		return errgo.Notef(err, "Unable to create request")
	}

	if p.options.Username != "" || p.options.Password != "" {
		req.SetBasicAuth(p.options.Username, p.options.Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return errgo.Notef(err, "Unable to send request")
	}
	defer resp.Body.Close()

	if p.options.ExpectedStatusCode == 0 {
		if resp.Status[0] != '2' && resp.Status[0] != '3' {
			return errgo.Newf("Invalid return code: %s", resp.Status)
		}
	} else {
		if resp.StatusCode != p.options.ExpectedStatusCode {
			return errgo.Newf("Unexpected status code: %v (expected: %v)", resp.StatusCode, p.options.ExpectedStatusCode)
		}
	}

	if p.options.Checker != nil {
		err = p.options.Checker.Check(resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewTimeoutClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, 2*time.Second)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(1 * time.Second))
				return conn, nil
			},
		},
	}
}
