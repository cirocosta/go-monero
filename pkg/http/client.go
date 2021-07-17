package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type HTTPClientConfig struct {
	TLSSkipVerify  bool
	TLSClientCert  string
	TLSClientKey   string
	TLSCACert      string
	Verbose        bool
	RequestTimeout time.Duration
}

func (c HTTPClientConfig) Validate() error {
	if c.TLSClientCert != "" && c.TLSClientKey == "" {
		return fmt.Errorf("tls client certificate specified but tls client key not")
	}

	if c.TLSClientKey != "" && c.TLSClientCert == "" {
		return fmt.Errorf("tls client key specified but tls client key")
	}

	return nil
}

// NewHTTPClient instantiates a new `http.Client` with a few defaults.
//
// `verbose`: if set, adds a transport that dumps all requests and responses to
// stdout.
//
// `skipTLSVerify`: if set, skips TLS validation.
//
func NewHTTPClient(cfg HTTPClientConfig) (*http.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	tlsConfig := &tls.Config{}
	if cfg.TLSCACert != "" {
		if err := WithCACert(cfg.TLSCACert)(tlsConfig); err != nil {
			return nil, fmt.Errorf("with tls ca cert: %w", err)
		}
	}

	if cfg.TLSClientCert != "" {
		if err := WithClientCertificate(cfg.TLSClientCert, cfg.TLSClientKey)(tlsConfig); err != nil {
			return nil, fmt.Errorf("with tls client certificate: %w", err)
		}
	}

	if cfg.TLSSkipVerify {
		WithInsecureSkipVerify()(tlsConfig)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = tlsConfig

	client := &http.Client{
		Timeout:   cfg.RequestTimeout,
		Transport: transport,
	}

	if cfg.Verbose {
		client.Transport = NewDumpTransport(client.Transport)
	}

	return client, nil
}
