package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type TLSOption func(*tls.Config) error

func WithCACert(fpath string) TLSOption {
	return func(config *tls.Config) error {
		certBytes, err := os.ReadFile(fpath)
		if err != nil {
			return fmt.Errorf("read file '%s': %w", fpath, err)
		}

		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(certBytes); !ok {
			return fmt.Errorf("ca cert '%s' not valid", fpath)
		}

		config.RootCAs = pool

		return nil
	}
}

func WithInsecureSkipVerify() TLSOption {
	return func(config *tls.Config) error {
		config.InsecureSkipVerify = true

		return nil
	}
}

func WithClientCertificate(cert, key string) TLSOption {
	return func(config *tls.Config) error {
		keypair, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return fmt.Errorf("load x509 key pair: %w", err)
		}

		config.Certificates = []tls.Certificate{keypair}

		return nil
	}
}
