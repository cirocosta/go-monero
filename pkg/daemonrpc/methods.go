package daemonrpc

import "fmt"

const (
	MethodGetBlockCount    = "get_block_count"
	MethodOnGetBlockHash   = "on_get_block_hash"
	MethodGetBlockTemplate = "get_block_template"
)

type GetBlockCountResult struct {
	Count  uint64 `json:"count"`
	Status string `json:"status"`
}

func (c *Client) GetBlockCount() (*GetBlockCountResult, error) {
	var (
		resp = &GetBlockCountResult{}
	)

	if err := c.JsonRPC(MethodGetBlockCount, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

func (c *Client) OnGetBlockHash(height uint64) (string, error) {
	var (
		resp   = ""
		params = []uint64{height}
	)

	if err := c.JsonRPC(MethodOnGetBlockHash, params, &resp); err != nil {
		return "", fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetBlockTemplateResult struct {
	BlockhashingBlob  string `json:"blockhashing_blob"`
	BlocktemplateBlob string `json:"blocktemplate_blob"`
	Difficulty        int64  `json:"difficulty"`
	ExpectedReward    int64  `json:"expected_reward"`
	Height            int    `json:"height"`
	PrevHash          string `json:"prev_hash"`
	ReservedOffset    int    `json:"reserved_offset"`
	Status            string `json:"status"`
	Untrusted         bool   `json:"untrusted"`
}

func (c *Client) GetBlockTemplate(walletAddress string, reserveSize uint) (*GetBlockTemplateResult, error) {
	var (
		resp   = &GetBlockTemplateResult{}
		params = map[string]interface{}{
			"wallet_address": walletAddress,
			"reserve_size":   reserveSize,
		}
	)

	if err := c.JsonRPC(MethodGetBlockTemplate, params, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}
