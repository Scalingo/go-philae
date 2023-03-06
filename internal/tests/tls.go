package tests

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

var testSubject = pkix.Name{
	Organization:  []string{"Scaling Test"},
	Country:       []string{"FR"},
	Province:      []string{""},
	Locality:      []string{"Strasbourg"},
	StreetAddress: []string{"3 place de Haguenau"},
	PostalCode:    []string{"67000"},
}

type CertificateWithKey struct {
	Certificate    *x509.Certificate
	PrivateKey     *rsa.PrivateKey
	PrivateKeyPEM  []byte
	CertificatePEM []byte
}

func GenerateCA() (*CertificateWithKey, error) {
	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               testSubject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}
	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		return nil, err
	}
	return &CertificateWithKey{
		Certificate:    ca,
		PrivateKey:     caPrivKey,
		PrivateKeyPEM:  caPrivKeyPEM.Bytes(),
		CertificatePEM: caPEM.Bytes(),
	}, nil
}

func GenerateCertificateWithCA(ca *CertificateWithKey) (*CertificateWithKey, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      testSubject,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca.Certificate, &certPrivKey.PublicKey, ca.PrivateKey)
	if err != nil {
		return nil, err
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return nil, err
	}

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		return nil, err
	}

	return &CertificateWithKey{
		Certificate:    cert,
		PrivateKey:     certPrivKey,
		PrivateKeyPEM:  certPrivKeyPEM.Bytes(),
		CertificatePEM: certPEM.Bytes(),
	}, nil
}
