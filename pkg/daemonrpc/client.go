package daemonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	EndpointJsonRPC = "/json_rpc"
	VersionJsonRPC  = "2.0"

	MethodGetBlockCount    = "get_block_count"
	MethodOnGetBlockHash   = "on_get_block_hash"
	MethodGetBlockTemplate = "get_block_template"
)

type Client struct {
	http *http.Client
	url  *url.URL
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

func NewHTTPClient(verbose bool) *http.Client {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	if verbose == true {
		client.Transport = NewDumpTransport(http.DefaultTransport)
	}

	return client

}

func NewClient(address string, opts ...ClientOption) (*Client, error) {
	options := &ClientOptions{
		HTTPClient: NewHTTPClient(false),
	}

	for _, opt := range opts {
		opt(options)
	}

	parsedAddress, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("url parse: %w", err)
	}

	return &Client{
		url:  parsedAddress,
		http: options.HTTPClient,
	}, nil
}

// ResponseEnvelope wraps all responses from the RPC server.
//
type ResponseEnvelope struct {
	Id      string      `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
}

// RequestEnvelope wraps all requests made to the RPC server.
//
type RequestEnvelope struct {
	Id      string                 `json:"id"`
	JsonRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type GetBlockCountResponse struct {
	Count  uint64 `json:"count"`
	Status string `json:"status"`
}

func (c *Client) GetBlockCount() (*GetBlockCountResponse, error) {
	resp := &GetBlockCountResponse{}

	if err := c.JsonRPC(MethodGetBlockCount, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

func (c *Client) JsonRPC(method string, params map[string]interface{}, response interface{}) error {
	url := *c.url
	url.Path = EndpointJsonRPC

	b, err := json.Marshal(&RequestEnvelope{
		Id:     "0",
		Method: method,
		Params: params,
	})
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest("GET", url.String(), bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("new req '%s': %w", url.String(), err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("non-2xx status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&ResponseEnvelope{
		Result: response,
	}); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
