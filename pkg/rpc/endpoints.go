package rpc

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	endpointGetHeight               = "/get_height"
	endpointGetPeerList             = "/get_peer_list"
	endpointGetTransactionPool      = "/get_transaction_pool"
	endpointGetTransactionPoolStats = "/get_transaction_pool_stats"
	endpointGetTransactions         = "/get_transactions"
	endpointGetNetStats             = "/get_net_stats"
)

type GetTransactionPoolResult struct {
	Credits        int `json:"credits"`
	SpentKeyImages []struct {
		IDHash    string   `json:"id_hash"`
		TxsHashes []string `json:"txs_hashes"`
	} `json:"spent_key_images"`
	Status       string `json:"status"`
	TopHash      string `json:"top_hash"`
	Transactions []struct {
		BlobSize           int    `json:"blob_size"`
		DoNotRelay         bool   `json:"do_not_relay"`
		DoubleSpendSeen    bool   `json:"double_spend_seen"`
		Fee                int    `json:"fee"`
		IDHash             string `json:"id_hash"`
		KeptByBlock        bool   `json:"kept_by_block"`
		LastFailedHeight   int    `json:"last_failed_height"`
		LastFailedIDHash   string `json:"last_failed_id_hash"`
		LastRelayedTime    int    `json:"last_relayed_time"`
		MaxUsedBlockHeight int    `json:"max_used_block_height"`
		MaxUsedBlockIDHash string `json:"max_used_block_id_hash"`
		ReceiveTime        int    `json:"receive_time"`
		Relayed            bool   `json:"relayed"`
		TxBlob             string `json:"tx_blob"`
		TxJSON             string `json:"tx_json"`
		Weight             int    `json:"weight"`
	} `json:"transactions"`
	Untrusted bool `json:"untrusted"`
}

func (c *Client) GetTransactionPool(ctx context.Context) (*GetTransactionPoolResult, error) {
	resp := &GetTransactionPoolResult{}

	if err := c.Request(ctx, endpointGetTransactionPool, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

type GetTransactionPoolStatsResult struct {
	Credits   int `json:"credits"`
	PoolStats struct {
		BytesMax   int `json:"bytes_max"`
		BytesMed   int `json:"bytes_med"`
		BytesMin   int `json:"bytes_min"`
		BytesTotal int `json:"bytes_total"`
		FeeTotal   int `json:"fee_total"`
		Histo      []struct {
			Bytes int `json:"bytes"`
			Txs   int `json:"txs"`
		} `json:"histo"`
		Histo98Pc       int `json:"histo_98pc"`
		Num10M          int `json:"num_10m"`
		NumDoubleSpends int `json:"num_double_spends"`
		NumFailing      int `json:"num_failing"`
		NumNotRelayed   int `json:"num_not_relayed"`
		Oldest          int `json:"oldest"`
		TxsTotal        int `json:"txs_total"`
	} `json:"pool_stats"`
	Status    string `json:"status"`
	TopHash   string `json:"top_hash"`
	Untrusted bool   `json:"untrusted"`
}

func (c *Client) GetTransactionPoolStats(ctx context.Context) (*GetTransactionPoolStatsResult, error) {
	resp := new(GetTransactionPoolStatsResult)

	if err := c.Request(ctx, endpointGetTransactionPoolStats, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

type GetPeerListResult struct {
	GrayList []struct {
		Host     string `json:"host"`
		ID       uint64 `json:"id"`
		IP       int    `json:"ip"`
		LastSeen int    `json:"last_seen"`
		Port     int    `json:"port"`
	} `json:"gray_list"`
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
	WhiteList []struct {
		Host        string `json:"host"`
		ID          uint64 `json:"id"`
		IP          int    `json:"ip"`
		LastSeen    int    `json:"last_seen"`
		Port        int    `json:"port"`
		PruningSeed int    `json:"pruning_seed"`
		RPCPort     int    `json:"rpc_port"`
	} `json:"white_list"`
}

func (c *Client) GetPeerList(ctx context.Context) (*GetPeerListResult, error) {
	resp := &GetPeerListResult{}

	if err := c.Request(ctx, endpointGetPeerList, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

type GetHeightResult struct {
	Hash      string `json:"hash"`
	Height    uint64 `json:"height"`
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
}

func (c *Client) GetHeight(ctx context.Context) (*GetHeightResult, error) {
	resp := &GetHeightResult{}

	if err := c.Request(ctx, endpointGetHeight, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

type GetNetStatsResult struct {
	StartTime       int    `json:"start_time"`
	Status          string `json:"status"`
	TotalBytesIn    uint64 `json:"total_bytes_in"`
	TotalBytesOut   uint64 `json:"total_bytes_out"`
	TotalPacketsIn  uint64 `json:"total_packets_in"`
	TotalPacketsOut uint64 `json:"total_packets_out"`
	Untrusted       bool   `json:"untrusted"`
}

func (c *Client) GetNetStats(ctx context.Context) (*GetNetStatsResult, error) {
	resp := &GetNetStatsResult{}

	if err := c.Request(ctx, endpointGetNetStats, nil, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}

type GetTransactionsResult struct {
	Credits int    `json:"credits"`
	Status  string `json:"status"`
	TopHash string `json:"top_hash"`
	Txs     []struct {
		AsHex           string `json:"as_hex"`
		AsJSON          string `json:"as_json"`
		BlockHeight     int    `json:"block_height"`
		BlockTimestamp  int    `json:"block_timestamp"`
		DoubleSpendSeen bool   `json:"double_spend_seen"`
		InPool          bool   `json:"in_pool"`
		OutputIndices   []int  `json:"output_indices"`
		PrunableAsHex   string `json:"prunable_as_hex"`
		PrunableHash    string `json:"prunable_hash"`
		PrunedAsHex     string `json:"pruned_as_hex"`
		TxHash          string `json:"tx_hash"`
	} `json:"txs"`
	TxsAsHex  []string `json:"txs_as_hex"`
	Untrusted bool     `json:"untrusted"`
}

func (r *GetTransactionsResult) GetTransactions() ([]*GetTransactionsResultJSONTxn, error) {
	txns := make([]*GetTransactionsResultJSONTxn, len(r.Txs))

	for idx, txn := range r.Txs {
		if len(txn.AsJSON) == 0 {
			return nil, fmt.Errorf("txn '%s' w/ empty `.as_json`", txn.TxHash)
		}

		t := &GetTransactionsResultJSONTxn{}
		if err := json.Unmarshal([]byte(txn.AsJSON), t); err != nil {
			return nil, fmt.Errorf("unmarshal txn '%s': %w", txn.TxHash, err)
		}

		txns[idx] = t
	}

	return txns, nil
}

type GetTransactionsResultJSONTxn struct {
	Version    int `json:"version"`
	UnlockTime int `json:"unlock_time"`
	Vin        []struct {
		Key struct {
			Amount     int    `json:"amount"`
			KeyOffsets []int  `json:"key_offsets"`
			KImage     string `json:"k_image"`
		} `json:"key"`
	} `json:"vin"`
	Vout []struct {
		Amount int `json:"amount"`
		Target struct {
			Key string `json:"key"`
		} `json:"target"`
	} `json:"vout"`
	Extra         []int `json:"extra"`
	RctSignatures struct {
		Type     int `json:"type"`
		Txnfee   int `json:"txnFee"`
		Ecdhinfo []struct {
			Amount string `json:"amount"`
		} `json:"ecdhInfo"`
		Outpk []string `json:"outPk"`
	} `json:"rct_signatures"`
	RctsigPrunable struct {
		Nbp int `json:"nbp"`
		Bp  []struct {
			A      string   `json:"A"`
			S      string   `json:"S"`
			T1     string   `json:"T1"`
			T2     string   `json:"T2"`
			Taux   string   `json:"taux"`
			Mu     string   `json:"mu"`
			L      []string `json:"L"`
			R      []string `json:"R"`
			LowerA string   `json:"a"`
			B      string   `json:"b"`
			T      string   `json:"t"`
		} `json:"bp"`
		Clsags []struct {
			S  []string `json:"s"`
			C1 string   `json:"c1"`
			D  string   `json:"D"`
		} `json:"CLSAGs"`
		Pseudoouts []string `json:"pseudoOuts"`
	} `json:"rctsig_prunable"`
}

func (c *Client) GetTransactions(ctx context.Context, txns []string) (*GetTransactionsResult, error) {
	resp := &GetTransactionsResult{}

	if err := c.Request(ctx, endpointGetTransactions, map[string]interface{}{
		"txs_hashes":     txns,
		"decode_as_json": true,
	}, resp); err != nil {
		return nil, fmt.Errorf("other: %w", err)
	}

	return resp, nil
}
