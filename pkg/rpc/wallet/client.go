package wallet

import "context"

// Requester is responsible for making HTTP requests to `monero-wallet-rpc`
// JSONRPC endpoints.
//
type Requester interface {
	// JSONRPC is used for callind methods under `/json_rpc` that follow
	// monero's `v2` response and error encapsulation.
	//
	JSONRPC(
		ctx context.Context, method string, params, result interface{},
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
