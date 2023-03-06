package tests

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"

	"github.com/pkg/errors"
)

func NewUnstartedServerWithTLSConfig(handler http.Handler) (*CertificateWithKey, *httptest.Server, error) {
	ca, err := GenerateCA()
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to generate CA")
	}

	cert, err := GenerateCertificateWithCA(ca)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to generate server certificate")
	}

	serverCert, err := tls.X509KeyPair(cert.CertificatePEM, cert.PrivateKeyPEM)
	if err != nil {
		return nil, nil, errors.Wrap(err, "fail to generate server keypair")
	}

	serverTLSConf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
	}

	serv := httptest.NewUnstartedServer(handler)
	serv.TLS = serverTLSConf
	return ca, serv, nil
}
