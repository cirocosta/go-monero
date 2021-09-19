package daemon

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	endpointGetHeight               = "/get_height"
	endpointGetLimit                = "/get_limit"
	endpointGetNetStats             = "/get_net_stats"
	endpointGetOuts                 = "/get_outs"
	endpointGetPeerList             = "/get_peer_list"
	endpointGetPublicNodes          = "/get_public_nodes"
	endpointGetTransactionPool      = "/get_transaction_pool"
	endpointGetTransactionPoolStats = "/get_transaction_pool_stats"
	endpointGetTransactions         = "/get_transactions"
	endpointMiningStatus            = "/mining_status"
	endpointSetLimit                = "/set_limit"
	endpointSetLogLevel             = "/set_log_level"
	endpointSetLogCategories        = "/set_log_categories"
	endpointStartMining             = "/start_mining"
	endpointStopMining              = "/stop_mining"
)

func (c *Client) StopMining(
	ctx context.Context,
) (*StopMiningResult, error) {
	resp := &StopMiningResult{}

	err := c.RawRequest(ctx, endpointStopMining, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetLimit(ctx context.Context) (*GetLimitResult, error) {
	resp := &GetLimitResult{}

	err := c.RawRequest(ctx, endpointGetLimit, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) SetLogCategories(
	ctx context.Context, params SetLogCategoriesRequestParameters,
) (*SetLogCategoriesResult, error) {
	resp := &SetLogCategoriesResult{}

	err := c.RawRequest(ctx, endpointSetLogCategories, params, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) SetLogLevel(
	ctx context.Context, params SetLogLevelRequestParameters,
) (*SetLogLevelResult, error) {
	resp := &SetLogLevelResult{}

	err := c.RawRequest(ctx, endpointSetLogLevel, params, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) SetLimit(
	ctx context.Context, params SetLimitRequestParameters,
) (*SetLimitResult, error) {
	resp := &SetLimitResult{}

	err := c.RawRequest(ctx, endpointSetLimit, params, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) StartMining(
	ctx context.Context, params StartMiningRequestParameters,
) (*StartMiningResult, error) {
	resp := &StartMiningResult{}

	err := c.RawRequest(ctx, endpointStartMining, params, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) MiningStatus(
	ctx context.Context,
) (*MiningStatusResult, error) {
	resp := &MiningStatusResult{}

	err := c.RawRequest(ctx, endpointMiningStatus, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetTransactionPool(
	ctx context.Context,
) (*GetTransactionPoolResult, error) {
	resp := &GetTransactionPoolResult{}

	err := c.RawRequest(ctx, endpointGetTransactionPool, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetTransactionPoolStats(
	ctx context.Context,
) (*GetTransactionPoolStatsResult, error) {
	resp := &GetTransactionPoolStatsResult{}

	err := c.RawRequest(ctx, endpointGetTransactionPoolStats, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetPeerList(
	ctx context.Context,
) (*GetPeerListResult, error) {
	resp := &GetPeerListResult{}

	err := c.RawRequest(ctx, endpointGetPeerList, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

type GetPublicNodesRequestParameters struct {
	Gray           bool `json:"gray"`
	White          bool `json:"white"`
	IncludeBlocked bool `json:"include_blocked"`
}

func (c *Client) GetPublicNodes(
	ctx context.Context, params GetPublicNodesRequestParameters,
) (*GetPublicNodesResult, error) {
	resp := &GetPublicNodesResult{}

	err := c.RawRequest(ctx, endpointGetPublicNodes, params, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetOuts(
	ctx context.Context, outputs []uint, gettxid bool,
) (*GetOutsResult, error) {
	resp := &GetOutsResult{}

	type output struct {
		Index uint `json:"index"`
	}

	params := struct {
		Outputs []output `json:"outputs"`
		GetTxID bool     `json:"get_txid,omitempty"`
	}{GetTxID: gettxid}

	for _, out := range outputs {
		params.Outputs = append(params.Outputs, output{out})
	}

	err := c.RawRequest(ctx, endpointGetOuts, params, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetHeight(ctx context.Context) (*GetHeightResult, error) {
	resp := &GetHeightResult{}

	err := c.RawRequest(ctx, endpointGetHeight, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (c *Client) GetNetStats(ctx context.Context) (*GetNetStatsResult, error) {
	resp := &GetNetStatsResult{}

	err := c.RawRequest(ctx, endpointGetNetStats, nil, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}

func (r *GetTransactionsResult) GetTransactions() ([]*TransactionJSON, error) {
	txns := make([]*TransactionJSON, len(r.Txs))

	for idx, txn := range r.Txs {
		if len(txn.AsJSON) == 0 {
			return nil, fmt.Errorf("txn w/ empty `.as_json`: %s",
				txn.TxHash)
		}

		t := &TransactionJSON{}
		err := json.Unmarshal([]byte(txn.AsJSON), t)
		if err != nil {
			return nil, fmt.Errorf("unmarshal txn '%s': %w",
				txn.TxHash, err)
		}

		txns[idx] = t
	}

	return txns, nil
}

func (c *Client) GetTransactions(
	ctx context.Context, txns []string,
) (*GetTransactionsResult, error) {
	resp := &GetTransactionsResult{}
	params := map[string]interface{}{
		"txs_hashes":     txns,
		"decode_as_json": true,
	}

	err := c.RawRequest(ctx, endpointGetTransactions, params, resp)
	if err != nil {
		return nil, fmt.Errorf("raw request: %w", err)
	}

	return resp, nil
}
