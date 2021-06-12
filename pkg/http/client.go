package http

import (
	"net/http"
	"time"
)

func NewHTTPClient(verbose bool) *http.Client {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	if verbose {
		client.Transport = NewDumpTransport(http.DefaultTransport)
	}

	return client
}
