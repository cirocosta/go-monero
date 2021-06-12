package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	mhttp "github.com/cirocosta/go-monero/pkg/http"
)

const (
	// endpointJsonRPC is the common endpoint used for all the RPC calls
	// that make use of epee's JSONRPC invocation format for requests and
	// responses.
	//
	endpointJsonRPC = "/json_rpc"

	// versionJsonRPC is the version of the JsonRPC format.
	//
	versionJsonRPC = "2.0"
)

// Client is a wrapper over a plain HTTP client providing methods that
// correspond to all RPC invocations to a `monerod` daemon, including
// restricted and non-restricted ones.
//
type Client struct {
	// http is the underlying http client that takes care of sending
	// requests and receiving the responses.
	//
	// To provide your own, make use of `WithHTTPClient` when instantiating
	// the client via the `NewClient` constructor.
	//
	http *http.Client

	// address is the address of the monerod instance serving the RPC
	// endpoints.
	//
	address *url.URL
}

type ClientOptions struct {
	HTTPClient *http.Client
}

type ClientOption func(o *ClientOptions)

func WithHTTPClient(v *http.Client) func(o *ClientOptions) {
	return func(o *ClientOptions) {
		o.HTTPClient = v
	}
}

// NewClient instantiates a new Client that is able to communicate with
// monerod's RPC endpoints.
//
// The `address` might be either restricted (typically <ip>:18089) or not
// (typically <ip>:18081).
//
func NewClient(address string, opts ...ClientOption) (*Client, error) {
	options := &ClientOptions{
		HTTPClient: mhttp.NewHTTPClient(false),
	}

	for _, opt := range opts {
		opt(options)
	}

	parsedAddress, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("url parse: %w", err)
	}

	return &Client{
		address: parsedAddress,
		http:    options.HTTPClient,
	}, nil
}

// ResponseEnvelope wraps all responses from the RPC server.
//
type ResponseEnvelope struct {
	Id      string      `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// RequestEnvelope wraps all requests made to the RPC server.
//
type RequestEnvelope struct {
	Id      string      `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// Other makes requests to any other endpoints that are not `/jsonrpc`.
//
func (c *Client) Other(ctx context.Context, endpoint string, params interface{}, response interface{}) error {
	address := *c.address
	address.Path = endpoint

	var body io.Reader

	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}

		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", address.String(), body)
	if err != nil {
		return fmt.Errorf("new req '%s': %w", address.String(), err)
	}

	req.Header.Add("Content-Type", "application/json")

	if err := c.submitRequest(req, response); err != nil {
		return fmt.Errorf("submit request: %w", err)
	}

	return nil
}

func (c *Client) JsonRPC(ctx context.Context, method string, params interface{}, response interface{}) error {
	address := *c.address
	address.Path = endpointJsonRPC

	b, err := json.Marshal(&RequestEnvelope{
		Id:      "0",
		JsonRPC: versionJsonRPC,
		Method:  method,
		Params:  params,
	})
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", address.String(), bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("new req '%s': %w", address.String(), err)
	}

	req.Header.Add("Content-Type", "application/json")

	rpcResponseBody := &ResponseEnvelope{
		Result: response,
	}

	if err := c.submitRequest(req, rpcResponseBody); err != nil {
		return fmt.Errorf("submit request: %w", err)
	}

	if rpcResponseBody.Error.Code != 0 || rpcResponseBody.Error.Message != "" {
		return fmt.Errorf("rpc error: code=%d message=%s",
			rpcResponseBody.Error.Code,
			rpcResponseBody.Error.Message,
		)
	}

	return nil
}

func (c *Client) submitRequest(req *http.Request, response interface{}) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("non-2xx status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
