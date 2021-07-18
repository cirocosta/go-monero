package daemon

import "context"

// Requester is responsible for making concrete request to Monero's endpoints,
// i.e., either `jsonrpc` methods or those "raw" endpoints.
//
type Requester interface {
	// JSONRPC is used for callind methods under `/json_rpc` that follow
	// monero's `v2` response and error encapsulation.
	//
	JSONRPC(
		ctx context.Context, method string, params, result interface{},
	) error

	// RawRequest is used for making a request to an arbitrary endpoint
	// `endpoint` whose response (in JSON format) should be unmarshalled to
	// `response`.
	//
	RawRequest(
		ctx context.Context,
		endpoint string,
		params interface{},
		response interface{},
	) error
}

// Client provides access to the daemon's JSONRPC methods and regular
// endpoints.
//
type Client struct {
	Requester
}

// NewClient instantiates a new client for interacting with monero's daemon
// api.
//
func NewClient(c Requester) *Client {
	return &Client{
		Requester: c,
	}
}
