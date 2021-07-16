package daemon

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	endpointGetHeight               = "/get_height"
	endpointGetOuts                 = "/get_outs"
	endpointGetNetStats             = "/get_net_stats"
	endpointGetPeerList             = "/get_peer_list"
	endpointGetPublicNodes          = "/get_public_nodes"
	endpointGetTransactionPool      = "/get_transaction_pool"
	endpointGetTransactionPoolStats = "/get_transaction_pool_stats"
	endpointGetTransactions         = "/get_transactions"
)

func (c *Client) GetTransactionPool(ctx context.Context) (*GetTransactionPoolResult, error) {
	resp := &GetTransactionPoolResult{}

	if err := c.RawRequest(ctx, endpointGetTransactionPool, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

func (c *Client) GetTransactionPoolStats(ctx context.Context) (*GetTransactionPoolStatsResult, error) {
	resp := new(GetTransactionPoolStatsResult)

	if err := c.RawRequest(ctx, endpointGetTransactionPoolStats, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

func (c *Client) GetPeerList(ctx context.Context) (*GetPeerListResult, error) {
	resp := &GetPeerListResult{}

	if err := c.RawRequest(ctx, endpointGetPeerList, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

type GetPublicNodesRequestParameters struct {
	Gray           bool `json:"gray"`
	White          bool `json:"white"`
	IncludeBlocked bool `json:"include_blocked"`
}

func (c *Client) GetPublicNodes(ctx context.Context, params GetPublicNodesRequestParameters) (*GetPublicNodesResult, error) {
	resp := &GetPublicNodesResult{}

	if err := c.RawRequest(ctx, endpointGetPublicNodes, params, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

func (c *Client) GetOuts(ctx context.Context, outputs []uint, gettxid bool) (*GetOutsResult, error) {
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

	if err := c.RawRequest(ctx, endpointGetOuts, params, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

func (c *Client) GetHeight(ctx context.Context) (*GetHeightResult, error) {
	resp := &GetHeightResult{}

	if err := c.RawRequest(ctx, endpointGetHeight, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

func (c *Client) GetNetStats(ctx context.Context) (*GetNetStatsResult, error) {
	resp := &GetNetStatsResult{}

	if err := c.RawRequest(ctx, endpointGetNetStats, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

func (r *GetTransactionsResult) GetTransactions() ([]*TransactionJSON, error) {
	txns := make([]*TransactionJSON, len(r.Txs))

	for idx, txn := range r.Txs {
		if len(txn.AsJSON) == 0 {
			return nil, fmt.Errorf("txn '%s' w/ empty `.as_json`", txn.TxHash)
		}

		t := &TransactionJSON{}
		if err := json.Unmarshal([]byte(txn.AsJSON), t); err != nil {
			return nil, fmt.Errorf("unmarshal txn '%s': %w", txn.TxHash, err)
		}

		txns[idx] = t
	}

	return txns, nil
}

func (c *Client) GetTransactions(ctx context.Context, txns []string) (*GetTransactionsResult, error) {
	resp := &GetTransactionsResult{}

	if err := c.RawRequest(ctx, endpointGetTransactions, map[string]interface{}{
		"txs_hashes":     txns,
		"decode_as_json": true,
	}, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}
