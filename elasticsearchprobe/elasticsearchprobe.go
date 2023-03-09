package elasticsearchprobe

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/opensearch-project/opensearch-go"
	"github.com/pkg/errors"
)

type ElasticsearchProbe struct {
	name     string
	url      string
	caCert   []byte
	insecure bool
	certPool CertPoolGetter
}

type ProbeOpts func(*ElasticsearchProbe)

func WithInsecureSkipVerify() ProbeOpts {
	return func(esProbe *ElasticsearchProbe) {
		esProbe.insecure = true
	}
}

func WithCA(caCert []byte) ProbeOpts {
	return func(esProbe *ElasticsearchProbe) {
		esProbe.caCert = caCert
	}
}

// NewElasticsearchProbe instantiate a new elasticsearch probe:
// - name: probe name
// - url : connection string with the form "http://username:password@example.com"
// - opts: optionnal parameters such as providing a TLS CA certificate or allowing insecure connections
func NewElasticsearchProbe(name, url string, opts ...ProbeOpts) ElasticsearchProbe {
	esProbe := ElasticsearchProbe{
		name:     name,
		url:      url,
		insecure: false,
		caCert:   []byte(""),
		certPool: DefaultCertPoolGetter{},
	}
	for _, opt := range opts {
		opt(&esProbe)
	}

	return esProbe
}

func (p ElasticsearchProbe) Name() string {
	return p.name
}

func (p ElasticsearchProbe) Check(_ context.Context) error {
	var certPool *x509.CertPool
	if p.caCert != nil && len(p.caCert) != 0 {
		certPool = p.certPool.FromCustomCA(p.caCert)
	} else {
		var err error
		certPool, err = p.certPool.SystemPool()
		if err != nil {
			return errors.Wrap(err, "fail to use system certificate pool")
		}
	}

	cfg := opensearch.Config{
		Addresses: []string{p.url},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: p.insecure,
				RootCAs:            certPool,
			},
		},
	}

	osClient, err := opensearch.NewClient(cfg)
	if err != nil {
		return errors.Wrap(err, "fail to open a new connection to Elasticsearch")
	}

	_, err = osClient.Info()
	if err != nil {
		return errors.Wrap(err, "fail to get elasticsearch info")
	}
	return nil
}
