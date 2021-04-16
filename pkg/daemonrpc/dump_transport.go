package daemonrpc

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type DumpTransport struct {
	R http.RoundTripper
}

func NewDumpTransport(rt http.RoundTripper) *DumpTransport {
	return &DumpTransport{
		R: rt,
	}
}

func (d *DumpTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	dump, _ := httputil.DumpRequestOut(h, true)
	fmt.Println(string(dump))

	resp, err := d.R.RoundTrip(h)
	fmt.Println(resp.StatusCode)

	dump, _ = httputil.DumpResponse(resp, true)
	fmt.Println(string(dump))

	return resp, err
}
