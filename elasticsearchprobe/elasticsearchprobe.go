package elasticsearchprobe

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/opensearch-project/opensearch-go"
	"github.com/pkg/errors"
)

type ElasticsearchProbe struct {
	name      string
	url       string
	caCert    []byte
	insecure  bool
	certPool  CertPoolGetter
	client    *opensearch.Client
	clientErr error
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

func WithCertPoolGetter(certPool CertPoolGetter) ProbeOpts {
	return func(esProbe *ElasticsearchProbe) {
		esProbe.certPool = certPool
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

	esProbe.client, esProbe.clientErr = esProbe.createClient()

	return esProbe
}

func (p ElasticsearchProbe) Name() string {
	return p.name
}

func (p *ElasticsearchProbe) createClient() (*opensearch.Client, error) {
	var certPool *x509.CertPool
	if p.caCert != nil && len(p.caCert) != 0 {
		certPool = p.certPool.FromCustomCA(p.caCert)
	} else {
		var err error
		certPool, err = p.certPool.SystemPool()
		if err != nil {
			return nil, errors.Wrap(err, "fail to use system certificate pool")
		}
	}

	cfg := opensearch.Config{
		Addresses: []string{p.url},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: p.insecure,
				RootCAs:            certPool,
			},
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 2,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return opensearch.NewClient(cfg)
}

func (p *ElasticsearchProbe) Check(_ context.Context) error {
	if p.clientErr != nil {
		return errors.Wrap(p.clientErr, "fail to open a new connection to Elasticsearch")
	}

	_, err := p.client.Info()
	if err != nil {
		return errors.Wrap(err, "fail to get elasticsearch info")
	}
	return nil
}
