package elasticsearchprobe

import "crypto/x509"

type CertPoolGetter interface {
	FromCustomCA(caCert []byte) *x509.CertPool
	SystemPool() (*x509.CertPool, error)
}

type DefaultCertPoolGetter struct{}

func (DefaultCertPoolGetter) FromCustomCA(caCert []byte) *x509.CertPool {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)
	return certPool
}

func (DefaultCertPoolGetter) SystemPool() (*x509.CertPool, error) {
	return x509.SystemCertPool()
}
