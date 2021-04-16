package daemonrpc

import "fmt"

const (
	MethodGetBlockCount    = "get_block_count"
	MethodOnGetBlockHash   = "on_get_block_hash"
	MethodGetBlockTemplate = "get_block_template"
)

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

func (c *Client) OnGetBlockHash(height uint64) (string, error) {
	var resp string

	if err := c.JsonRPC(MethodOnGetBlockHash, []uint64{height}, &resp); err != nil {
		return "", fmt.Errorf("get: %w", err)
	}

	return resp, nil
}
