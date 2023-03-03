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
	url      string
	insecure bool
	caCert   []byte
}

func NewPinger(url string, insecure bool, caCert []byte) Pinger {
	return pinger{
		url:      url,
		insecure: insecure,
		caCert:   caCert,
	}
}

func (pg pinger) Ping() error {
	var certPool *x509.CertPool
	if pg.caCert != nil {
		certPool = x509.NewCertPool()
		certPool.AppendCertsFromPEM(pg.caCert)
	} else {
		var err error
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return errors.Wrap(err, "fail to use system certificate pool")
		}
	}

	cfg := opensearch.Config{
		Addresses: []string{pg.url},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: pg.insecure,
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
	}
	for _, opt := range opts {
		opt(&esProbe)
	}

	esPinger := NewPinger(esProbe.url, esProbe.insecure, esProbe.caCert)
	esProbe.pinger = esPinger
	return esProbe
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
