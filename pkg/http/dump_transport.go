package http

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

// DumpTransport implements the `net/http.RoundTripper` interface wrapping
// another Roundtripper dumping to stdout both the requests and the responses
// that it sees passing through.
//
type DumpTransport struct {
	R http.RoundTripper
}

// NewDumpTransport instantiates a new DumpTransport.
//
func NewDumpTransport(rt http.RoundTripper) *DumpTransport {
	return &DumpTransport{
		R: rt,
	}
}

// RoundTrip implements the functionality of dumping http requests and
// responses to `stdout` for each HTTP transaction that passes through it.
//
// It does so by first dumping the request, then passing that down to the
// wrapped roundtripper, and then from the response it sees, dumping it too.
//

func (d *DumpTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	requestDump, _ := httputil.DumpRequestOut(h, true)
	fmt.Fprintln(os.Stderr, string(requestDump))

	resp, err := d.R.RoundTrip(h)
	if err != nil {
		return nil, err
	}

	responseDump, _ := httputil.DumpResponse(resp, true)
	fmt.Fprintln(os.Stderr, string(responseDump))

	return resp, err
}
