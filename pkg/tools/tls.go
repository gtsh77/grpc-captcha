package tools

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
)

func PrepareTLS(keydata, crtdata, cadata, domain string) (*tls.Config, error) {
	var (
		crt    tls.Certificate
		ca     *x509.CertPool
		tlsCfg *tls.Config

		err error
	)

	ca = x509.NewCertPool()
	if crt, err = tls.X509KeyPair([]byte(crtdata), []byte(keydata)); err != nil {
		return nil, fmt.Errorf("tls.X509KeyPair: %w", err)
	}

	if ok := ca.AppendCertsFromPEM([]byte(cadata)); !ok {
		return nil, fmt.Errorf("x509.CertPool.AppendCertsFromPEM: %w", err)
	}

	tlsCfg = &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ServerName:   domain,
		RootCAs:      ca,
		Certificates: []tls.Certificate{crt},
	}

	return tlsCfg, nil
}
