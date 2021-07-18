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

	// Username is the name of the user to send in the header of every HTTP
	// call - must match the first portion of
	// `--rpc-login=<username>:[password]` provided to `monerod`.
	//
	Username string

	// Password is the user's password to send in the header of every HTTP
	// call - must match the second portion of
	// `--rpc-login=<username>:[password]` provided to `monerod` or the
	// password interactively supplied during the daemon's startup.
	//
	// Note that because the `monerod` performs digest auth, the password
	// won't be sent solely in plain base64 encoding, but the rest of the
	// body of every request and response will still be cleartext.
	//
	Password string
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

	if c.Username != "" && c.Password == "" {
		return fmt.Errorf("username specified but password not")
	}

	if c.Password != "" && c.Username == "" {
		return fmt.Errorf("password specified but username not")
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

	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
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

	if cfg.Username != "" {
		client.Transport = NewDigestAuthTransport(
			cfg.Username, cfg.Password,
			client.Transport,
		)
	}

	return client, nil
}

func WithTransport(rt http.RoundTripper) func(*http.Client) {
	return func(c *http.Client) {
		c.Transport = rt
	}
}
