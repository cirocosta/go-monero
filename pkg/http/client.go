package http

import (
	"net/http"
	"time"
)

// NewHTTPClient instantiates a new `http.Client` with a few defaults.
//
// `verbose`: if set, adds a transport that dumps all requests and responses to
// stdout.
//
func NewHTTPClient(verbose bool) *http.Client {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	if verbose {
		client.Transport = NewDumpTransport(http.DefaultTransport)
	}

	return client
}
