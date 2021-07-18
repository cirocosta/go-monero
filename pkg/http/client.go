package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// ClientConfig provides extra configuration to tweak the behavior of the HTTP
// client instantiated via `NewClient`.
//
type ClientConfig struct {
	// TLSSkipVerify indicates that the client should not perform any
	// hostname or certificate chain of trust validations.
	//
	TLSSkipVerify bool

	// TLSClientCert is the path to a TLS client certificate to be used
	// when connecting to a TLS server.
	//
	// ps.: must be supplied together with TLSClientKey.
	//
	TLSClientCert string

	// TLSClientKey is the path to a TLS private key certificate to be used
	// when connecting to a TLS server.
	//
	// ps.: must be supplied together with TLSClientCert.
	//
	TLSClientKey string

	// TLSCACert is the path to a certificate authority certificate that
	// should be included in the chain of trust.
	//
	TLSCACert string

	// Verbose dictates whether the transport should dump all request and
	// response information to stderr.
	//
	Verbose bool

	// RequestTimeout places a deadline on every request issued by this
	// client.
	//
	RequestTimeout time.Duration
}

func (c ClientConfig) Validate() error {
	if c.TLSClientCert != "" && c.TLSClientKey == "" {
		return fmt.Errorf("tls client certificate specified " +
			"but tls client key not")
	}

	if c.TLSClientKey != "" && c.TLSClientCert == "" {
		return fmt.Errorf("tls client key specified but " +
			"tls client key")
	}

	return nil
}

// NewClient instantiates a new `http.Client` based on the client configuration
// supplied.
//
func NewClient(cfg ClientConfig) (*http.Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if cfg.TLSCACert != "" {
		if err := WithCACert(cfg.TLSCACert)(tlsConfig); err != nil {
			return nil, fmt.Errorf("with tls ca cert: %w", err)
		}
	}

	if cfg.TLSClientCert != "" {
		err := WithClientCertificate(
			cfg.TLSClientCert, cfg.TLSClientKey,
		)(tlsConfig)
		if err != nil {
			return nil, fmt.Errorf(
				"with tls client certificate: %w", err)
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
