package daemon

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	methodGenerateBlocks         = "generateblocks"
	methodGetAlternateChains     = "get_alternate_chains"
	methodGetBans                = "get_bans"
	methodGetBlock               = "get_block"
	methodGetBlockCount          = "get_block_count"
	methodGetBlockHeaderByHash   = "get_block_header_by_hash"
	methodGetBlockHeaderByHeight = "get_block_header_by_height"
	methodGetBlockTemplate       = "get_block_template"
	methodGetCoinbaseTxSum       = "get_coinbase_tx_sum"
	methodGetConnections         = "get_connections"
	methodGetFeeEstimate         = "get_fee_estimate"
	methodGetInfo                = "get_info"
	methodGetLastBlockHeader     = "get_last_block_header"
	methodGetVersion             = "get_version"
	methodHardForkInfo           = "hard_fork_info"
	methodOnGetBlockHash         = "on_get_block_hash"
	methodRPCAccessTracking      = "rpc_access_tracking"
	methodRelayTx                = "relay_tx"
	methodSetBans                = "set_bans"
	methodSyncInfo               = "sync_info"
)

// GetAlternateChains displays alternative chains seen by the node.
//
// (restricted).
//
func (c *Client) GetAlternateChains(ctx context.Context) (*GetAlternateChainsResult, error) {
	resp := &GetAlternateChainsResult{}

	if err := c.JSONRPC(ctx, methodGetAlternateChains, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// RPCAccessTracking retrieves statistics that the monero daemon keeps track of
// about the use of each RPC method and endpoint.
//
// (restricted).
//
func (c *Client) RPCAccessTracking(ctx context.Context) (*RPCAccessTrackingResult, error) {
	resp := &RPCAccessTrackingResult{}

	if err := c.JSONRPC(ctx, methodRPCAccessTracking, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// HardForkInfo looks up informaiton about the last hard fork.
//
func (c *Client) HardForkInfo(ctx context.Context) (*HardForkInfoResult, error) {
	resp := &HardForkInfoResult{}

	if err := c.JSONRPC(ctx, methodHardForkInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetBans retrieves the list of banned IPs.
//
// (restricted).
//
func (c *Client) GetBans(ctx context.Context) (*GetBansResult, error) {
	resp := &GetBansResult{}

	if err := c.JSONRPC(ctx, methodGetBans, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

type SetBansBan struct {
	Host    string `json:"host"`
	Ban     bool   `json:"ban"`
	Seconds int64  `json:"seconds"`
}

type SetBansRequestParameters struct {
	Bans []SetBansBan `json:"bans"`
}

// SetBans bans a particular host.
//
// (restricted).
//
func (c *Client) SetBans(ctx context.Context, params SetBansRequestParameters) (*SetBansResult, error) {
	resp := &SetBansResult{}

	if err := c.JSONRPC(ctx, methodSetBans, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetVersion retrieves the version of monerod that the node uses.
//
// (restricted).
//
func (c *Client) GetVersion(ctx context.Context) (*GetVersionResult, error) {
	resp := &GetVersionResult{}

	if err := c.JSONRPC(ctx, methodGetVersion, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GenerateBlocksRequestParameters is the set of parameters to be passed to the
// GenerateBlocks RPC method.
//
type GenerateBlocksRequestParameters struct {
	// AmountOfBlocks is the number of blocks to be generated.
	//
	AmountOfBlocks uint64 `json:"amount_of_blocks,omitempty"`

	// WalletAddress is the address of the wallet that will get the rewards
	// of the coinbase transaction for such the blocks generates.
	//
	WalletAddress string `json:"wallet_address,omitempty"`

	// PreviousBlock TODO
	//
	PreviousBlock string `json:"prev_block,omitempty"`

	// StartingNonce TODO
	//
	StartingNonce uint32 `json:"starting_nonce,omitempty"`
}

// GenerateBlocks combines functionality from `GetBlockTemplate` and
// `SubmitBlock` RPC calls to allow rapid block creation.
//
// Difficulty is set permanently to 1 for regtest.
//
// (restricted).
//
func (c *Client) GenerateBlocks(ctx context.Context, params GenerateBlocksRequestParameters) (*GenerateBlocksResult, error) {
	resp := &GenerateBlocksResult{}

	if err := c.JSONRPC(ctx, methodGenerateBlocks, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetBlockCount(ctx context.Context) (*GetBlockCountResult, error) {
	resp := &GetBlockCountResult{}

	if err := c.JSONRPC(ctx, methodGetBlockCount, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) OnGetBlockHash(ctx context.Context, height uint64) (string, error) {
	var (
		resp   = ""
		params = []uint64{height}
	)

	if err := c.JSONRPC(ctx, methodOnGetBlockHash, params, &resp); err != nil {
		return "", fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) RelayTx(ctx context.Context, txns []string) (*RelayTxResult, error) {
	var (
		resp   = &RelayTxResult{}
		params = map[string]interface{}{
			"txids": txns,
		}
	)

	if err := c.JSONRPC(ctx, methodRelayTx, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetBlockTemplate gets a block template on which mining a new block.
//
func (c *Client) GetBlockTemplate(ctx context.Context, walletAddress string, reserveSize uint) (*GetBlockTemplateResult, error) {
	var (
		resp   = &GetBlockTemplateResult{}
		params = map[string]interface{}{
			"wallet_address": walletAddress,
			"reserve_size":   reserveSize,
		}
	)

	if err := c.JSONRPC(ctx, methodGetBlockTemplate, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetConnections(ctx context.Context) (*GetConnectionsResult, error) {
	resp := &GetConnectionsResult{}

	if err := c.JSONRPC(ctx, methodGetConnections, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetInfo retrieves general information about the state of the node and the
// network.
//
func (c *Client) GetInfo(ctx context.Context) (*GetInfoResult, error) {
	resp := &GetInfoResult{}

	if err := c.JSONRPC(ctx, methodGetInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetLastBlockHeader(ctx context.Context) (*GetLastBlockHeaderResult, error) {
	resp := &GetLastBlockHeaderResult{}

	if err := c.JSONRPC(ctx, methodGetLastBlockHeader, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetCoinbaseTxSum(ctx context.Context, height, count uint64) (*GetCoinbaseTxSumResult, error) {
	var (
		resp   = &GetCoinbaseTxSumResult{}
		params = map[string]uint64{
			"height": height,
			"count":  count,
		}
	)

	if err := c.JSONRPC(ctx, methodGetCoinbaseTxSum, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// InnerJSON parses the content of the JSON embedded in `GetBlockResult`.
//
func (j *GetBlockResult) InnerJSON() (*GetBlockResultJSON, error) {
	res := &GetBlockResultJSON{}

	if err := json.Unmarshal([]byte(j.JSON), res); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return res, nil
}

// GetBlockHeaderByHeight retrieves block header information for either one or
// multiple blocks.
//
func (c *Client) GetBlockHeaderByHeight(ctx context.Context, height uint64) (*GetBlockHeaderByHeightResult, error) {
	resp := &GetBlockHeaderByHeightResult{}

	if err := c.JSONRPC(ctx, methodGetBlockHeaderByHeight, map[string]interface{}{
		"height": height,
	}, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetBlockHeaderByHash retrieves block header information for either one or
// multiple blocks.
//
func (c *Client) GetBlockHeaderByHash(ctx context.Context, hashes []string) (*GetBlockHeaderByHashResult, error) {
	resp := &GetBlockHeaderByHashResult{}

	if err := c.JSONRPC(ctx, methodGetBlockHeaderByHash, map[string]interface{}{
		"hashes": hashes,
	}, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetBlockRequestParameters represents the set of possible parameters that can be used
// for submitting a call to the `get_block` jsonrpc method.
//
type GetBlockRequestParameters struct {
	Height uint64 `json:"height,omitempty"`
	Hash   string `json:"hash,omitempty"`
}

// GetBlock fetches full block information from a block at a particular hash OR height.
//
func (c *Client) GetBlock(ctx context.Context, params GetBlockRequestParameters) (*GetBlockResult, error) {
	resp := &GetBlockResult{}
	if err := c.JSONRPC(ctx, methodGetBlock, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetFeeEstimate(ctx context.Context, graceBlocks uint64) (*GetFeeEstimateResult, error) {
	var (
		resp   = new(GetFeeEstimateResult)
		params = map[string]uint64{
			"grace_blocks": graceBlocks,
		}
	)

	if err := c.JSONRPC(ctx, methodGetFeeEstimate, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) SyncInfo(ctx context.Context) (*SyncInfoResult, error) {
	resp := new(SyncInfoResult)

	if err := c.JSONRPC(ctx, methodSyncInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}
