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
func (c *Client) GetAlternateChains(
	ctx context.Context,
) (*GetAlternateChainsResult, error) {
	resp := &GetAlternateChainsResult{}

	err := c.JSONRPC(ctx, methodGetAlternateChains, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// RPCAccessTracking retrieves statistics that the monero daemon keeps track of
// about the use of each RPC method and endpoint.
//
// (restricted).
//
func (c *Client) RPCAccessTracking(
	ctx context.Context,
) (*RPCAccessTrackingResult, error) {
	resp := &RPCAccessTrackingResult{}

	err := c.JSONRPC(ctx, methodRPCAccessTracking, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// HardForkInfo looks up informaiton about the last hard fork.
//
func (c *Client) HardForkInfo(
	ctx context.Context,
) (*HardForkInfoResult, error) {
	resp := &HardForkInfoResult{}

	err := c.JSONRPC(ctx, methodHardForkInfo, nil, resp)
	if err != nil {
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

	err := c.JSONRPC(ctx, methodGetBans, nil, resp)
	if err != nil {
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
func (c *Client) SetBans(
	ctx context.Context, params SetBansRequestParameters,
) (*SetBansResult, error) {
	resp := &SetBansResult{}

	err := c.JSONRPC(ctx, methodSetBans, params, resp)
	if err != nil {
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

	err := c.JSONRPC(ctx, methodGetVersion, nil, resp)
	if err != nil {
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
func (c *Client) GenerateBlocks(
	ctx context.Context, params GenerateBlocksRequestParameters,
) (*GenerateBlocksResult, error) {
	resp := &GenerateBlocksResult{}

	err := c.JSONRPC(ctx, methodGenerateBlocks, params, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetBlockCount(
	ctx context.Context,
) (*GetBlockCountResult, error) {
	resp := &GetBlockCountResult{}

	err := c.JSONRPC(ctx, methodGetBlockCount, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) OnGetBlockHash(
	ctx context.Context, height uint64,
) (string, error) {
	resp := ""
	params := []uint64{height}

	err := c.JSONRPC(ctx, methodOnGetBlockHash, params, &resp)
	if err != nil {
		return "", fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) RelayTx(
	ctx context.Context, txns []string,
) (*RelayTxResult, error) {
	resp := &RelayTxResult{}
	params := map[string]interface{}{
		"txids": txns,
	}

	err := c.JSONRPC(ctx, methodRelayTx, params, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetBlockTemplate gets a block template on which mining a new block.
//
func (c *Client) GetBlockTemplate(
	ctx context.Context, walletAddress string, reserveSize uint,
) (*GetBlockTemplateResult, error) {
	resp := &GetBlockTemplateResult{}
	params := map[string]interface{}{
		"wallet_address": walletAddress,
		"reserve_size":   reserveSize,
	}

	err := c.JSONRPC(ctx, methodGetBlockTemplate, params, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetConnections(
	ctx context.Context,
) (*GetConnectionsResult, error) {
	resp := &GetConnectionsResult{}

	err := c.JSONRPC(ctx, methodGetConnections, nil, resp)
	if err != nil {
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

func (c *Client) GetLastBlockHeader(
	ctx context.Context,
) (*GetLastBlockHeaderResult, error) {
	resp := &GetLastBlockHeaderResult{}

	err := c.JSONRPC(ctx, methodGetLastBlockHeader, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetCoinbaseTxSum(
	ctx context.Context, height, count uint64,
) (*GetCoinbaseTxSumResult, error) {
	resp := &GetCoinbaseTxSumResult{}
	params := map[string]uint64{
		"height": height,
		"count":  count,
	}

	err := c.JSONRPC(ctx, methodGetCoinbaseTxSum, params, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// InnerJSON parses the content of the JSON embedded in `GetBlockResult`.
//
func (j *GetBlockResult) InnerJSON() (*GetBlockResultJSON, error) {
	res := &GetBlockResultJSON{}

	err := json.Unmarshal([]byte(j.JSON), res)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return res, nil
}

// GetBlockHeaderByHeight retrieves block header information for either one or
// multiple blocks.
//
func (c *Client) GetBlockHeaderByHeight(
	ctx context.Context, height uint64,
) (*GetBlockHeaderByHeightResult, error) {
	resp := &GetBlockHeaderByHeightResult{}
	params := map[string]interface{}{
		"height": height,
	}

	err := c.JSONRPC(ctx, methodGetBlockHeaderByHeight, params, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetBlockHeaderByHash retrieves block header information for either one or
// multiple blocks.
//
func (c *Client) GetBlockHeaderByHash(
	ctx context.Context, hashes []string,
) (*GetBlockHeaderByHashResult, error) {
	resp := &GetBlockHeaderByHashResult{}
	params := map[string]interface{}{
		"hashes": hashes,
	}

	err := c.JSONRPC(ctx, methodGetBlockHeaderByHash, params, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

// GetBlockRequestParameters represents the set of possible parameters that can
// be used for submitting a call to the `get_block` jsonrpc method.
//
type GetBlockRequestParameters struct {
	Height uint64 `json:"height,omitempty"`
	Hash   string `json:"hash,omitempty"`
}

// GetBlock fetches full block information from a block at a particular hash OR
// height.
//
func (c *Client) GetBlock(
	ctx context.Context, params GetBlockRequestParameters,
) (*GetBlockResult, error) {
	resp := &GetBlockResult{}

	err := c.JSONRPC(ctx, methodGetBlock, params, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) GetFeeEstimate(
	ctx context.Context, graceBlocks uint64,
) (*GetFeeEstimateResult, error) {
	resp := &GetFeeEstimateResult{}
	params := map[string]uint64{
		"grace_blocks": graceBlocks,
	}

	err := c.JSONRPC(ctx, methodGetFeeEstimate, params, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

func (c *Client) SyncInfo(ctx context.Context) (*SyncInfoResult, error) {
	resp := &SyncInfoResult{}

	err := c.JSONRPC(ctx, methodSyncInfo, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}
