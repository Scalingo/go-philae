package elasticsearchprobe

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/opensearch-project/opensearch-go"
	"github.com/pkg/errors"
)

type ElasticsearchProbe struct {
	name string
	url  string
}

// NewElasticsearchProbe instantiate a new elasticsearch probe:
// - name: probe name
// - url : connection string with the form "http://username:password@example.com"
func NewElasticsearchProbe(name, url string) ElasticsearchProbe {
	return ElasticsearchProbe{
		name: name,
		url:  url,
	}
}

func (p ElasticsearchProbe) Name() string {
	return p.name
}

func (p ElasticsearchProbe) Check(_ context.Context) error {
	cfg := opensearch.Config{
		Addresses: []string{p.url},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	osClient, err := opensearch.NewClient(cfg)
	if err != nil {
		return errors.Wrap(err, "fail to open a new connection to Elasticsearch")
	}
	_, err = osClient.Info()
	if err != nil {
		return errors.Wrap(err, "fail to get response")
	}
	return nil
}
