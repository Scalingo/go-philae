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
	pinger   Pinger
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

type Pinger interface {
	Ping() error
}
type pinger struct {
	opensearchClient opensearch.Client
}

func NewPinger(url string, insecure bool, certPool *x509.CertPool) (Pinger, error) {
	cfg := opensearch.Config{
		Addresses: []string{url},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
				RootCAs:            certPool,
			},
		},
	}
	osClient, err := opensearch.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "fail to open a new connection to Elasticsearch")
	}
	return pinger{opensearchClient: *osClient}, nil
}
func (p pinger) Ping() error {
	_, err := p.opensearchClient.Info()
	if err != nil {
		return errors.Wrap(err, "fail to get response")
	}
	return nil
}

// NewElasticsearchProbe instantiate a new elasticsearch probe:
// - name: probe name
// - url : connection string with the form "http://username:password@example.com"
// - opts: optionnal parameters such as providing a TLS CA certificate or allowing insecure connections
func NewElasticsearchProbe(name, url string, opts ...ProbeOpts) (ElasticsearchProbe, error) {
	esProbe := ElasticsearchProbe{
		name:     name,
		url:      url,
		insecure: false,
		caCert:   []byte(""),
	}
	for _, opt := range opts {
		opt(&esProbe)
	}

	var certPool *x509.CertPool
	if esProbe.caCert != nil {
		certPool = x509.NewCertPool()
		certPool.AppendCertsFromPEM(esProbe.caCert)
	} else {
		var err error
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return esProbe, errors.Wrap(err, "fail to use system certificate pool")
		}
	}
	esPinger, err := NewPinger(esProbe.url, esProbe.insecure, certPool)
	if err != nil {
		return esProbe, errors.Wrap(err, "fail to create elasticsearch client")
	}
	esProbe.pinger = esPinger
	return esProbe, nil
}

func (p ElasticsearchProbe) Name() string {
	return p.name
}

func (p ElasticsearchProbe) Check(_ context.Context) error {
	err := p.pinger.Ping()
	if err != nil {
		return errors.Wrap(err, "fail to get response")
	}
	return nil
}
