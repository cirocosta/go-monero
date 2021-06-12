package rpc

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	methodGetAlternateChains = "get_alternate_chains"
	methodGetBans            = "get_bans"
	methodGetBlock           = "get_block"
	methodGetBlockCount      = "get_block_count"
	methodGetBlockTemplate   = "get_block_template"
	methodGetCoinbaseTxSum   = "get_coinbase_tx_sum"
	methodGetConnections     = "get_connections"
	methodGetFeeEstimate     = "get_fee_estimate"
	methodGetInfo            = "get_info"
	methodGetLastBlockHeader = "get_last_block_header"
	methodHardForkInfo       = "hard_fork_info"
	methodOnGetBlockHash     = "on_get_block_hash"
	methodSyncInfo           = "sync_info"
	methodRPCAccessTracking  = "rpc_access_tracking"
	methodRelayTx            = "relay_tx"
)

type RPCAccessTrackingResult struct {
	Data []struct {
		Count   uint64 `json:"count"`
		Credits uint64 `json:"credits"`
		RPC     string `json:"rpc"`
		Time    uint64 `json:"time"`
	} `json:"data"`
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
}

func (c *Client) RPCAccessTracking(ctx context.Context) (*RPCAccessTrackingResult, error) {
	var resp = &RPCAccessTrackingResult{}

	if err := c.JsonRPC(ctx, methodRPCAccessTracking, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type HardForkInfoResult struct {
	Credits        int    `json:"credits"`
	EarliestHeight int    `json:"earliest_height"`
	Enabled        bool   `json:"enabled"`
	State          int    `json:"state"`
	Status         string `json:"status"`
	Threshold      int    `json:"threshold"`
	TopHash        string `json:"top_hash"`
	Untrusted      bool   `json:"untrusted"`
	Version        int    `json:"version"`
	Votes          int    `json:"votes"`
	Voting         int    `json:"voting"`
	Window         int    `json:"window"`
}

func (c *Client) HardForkInfo(ctx context.Context) (*HardForkInfoResult, error) {
	var resp = &HardForkInfoResult{}

	if err := c.JsonRPC(ctx, methodHardForkInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetBansResult struct {
	Bans []struct {
		Host    string `json:"host"`
		IP      int    `json:"ip"`
		Seconds int    `json:"seconds"`
	} `json:"bans"`
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
}

func (c *Client) GetBans(ctx context.Context) (*GetBansResult, error) {
	var resp = &GetBansResult{}

	if err := c.JsonRPC(ctx, methodGetBans, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetAlternateChainsResult struct {
	Chains []struct {
		BlockHash            string   `json:"block_hash"`
		BlockHashes          []string `json:"block_hashes"`
		Difficulty           int64    `json:"difficulty"`
		DifficultyTop64      int      `json:"difficulty_top64"`
		Height               int      `json:"height"`
		Length               int      `json:"length"`
		MainChainParentBlock string   `json:"main_chain_parent_block"`
		WideDifficulty       string   `json:"wide_difficulty"`
	} `json:"chains"`
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
}

func (c *Client) GetAlternateChains(ctx context.Context) (*GetAlternateChainsResult, error) {
	var resp = &GetAlternateChainsResult{}

	if err := c.JsonRPC(ctx, methodGetAlternateChains, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetBlockCountResult struct {
	Count  uint64 `json:"count"`
	Status string `json:"status"`
}

func (c *Client) GetBlockCount(ctx context.Context) (*GetBlockCountResult, error) {
	var resp = &GetBlockCountResult{}

	if err := c.JsonRPC(ctx, methodGetBlockCount, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

func (c *Client) OnGetBlockHash(ctx context.Context, height uint64) (string, error) {
	var (
		resp   = ""
		params = []uint64{height}
	)

	if err := c.JsonRPC(ctx, methodOnGetBlockHash, params, &resp); err != nil {
		return "", fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type RelayTxResult struct {
	Credits   int    `json:"credits"`
	Status    string `json:"status"`
	TopHash   string `json:"top_hash"`
	Untrusted bool   `json:"untrusted"`
}

func (c *Client) RelayTx(ctx context.Context, txns []string) (*RelayTxResult, error) {
	var (
		resp   = &RelayTxResult{}
		params = map[string]interface{}{
			"txids": txns,
		}
	)

	if err := c.JsonRPC(ctx, methodRelayTx, params, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
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

func (c *Client) GetBlockTemplate(ctx context.Context, walletAddress string, reserveSize uint) (*GetBlockTemplateResult, error) {
	var (
		resp   = &GetBlockTemplateResult{}
		params = map[string]interface{}{
			"wallet_address": walletAddress,
			"reserve_size":   reserveSize,
		}
	)

	if err := c.JsonRPC(ctx, methodGetBlockTemplate, params, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetConnectionsResult struct {
	Connections []struct {
		Address         string `json:"address"`
		AvgDownload     int    `json:"avg_download"`
		AvgUpload       int    `json:"avg_upload"`
		ConnectionID    string `json:"connection_id"`
		CurrentDownload int    `json:"current_download"`
		CurrentUpload   int    `json:"current_upload"`
		Height          int    `json:"height"`
		Host            string `json:"host"`
		Incoming        bool   `json:"incoming"`
		IP              string `json:"ip"`
		LiveTime        int    `json:"live_time"`
		LocalIP         bool   `json:"local_ip"`
		Localhost       bool   `json:"localhost"`
		PeerID          string `json:"peer_id"`
		Port            string `json:"port"`
		RecvCount       int    `json:"recv_count"`
		RecvIdleTime    int    `json:"recv_idle_time"`
		SendCount       int    `json:"send_count"`
		SendIdleTime    int    `json:"send_idle_time"`
		State           string `json:"state"`
		SupportFlags    int    `json:"support_flags"`
	} `json:"connections"`
	Status string `json:"status"`
}

func (c *Client) GetConnections(ctx context.Context) (*GetConnectionsResult, error) {
	var resp = &GetConnectionsResult{}

	if err := c.JsonRPC(ctx, methodGetConnections, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetInfoResult struct {
	AltBlocksCount           int    `json:"alt_blocks_count"`
	BlockSizeLimit           int    `json:"block_size_limit"`
	BlockSizeMedian          int    `json:"block_size_median"`
	BootstrapDaemonAddress   string `json:"bootstrap_daemon_address"`
	BusySyncing              bool   `json:"busy_syncing"`
	CumulativeDifficulty     int64  `json:"cumulative_difficulty"`
	Difficulty               int64  `json:"difficulty"`
	FreeSpace                int64  `json:"free_space"`
	GreyPeerlistSize         int    `json:"grey_peerlist_size"`
	Height                   int    `json:"height"`
	HeightWithoutBootstrap   int    `json:"height_without_bootstrap"`
	IncomingConnectionsCount int    `json:"incoming_connections_count"`
	Mainnet                  bool   `json:"mainnet"`
	Offline                  bool   `json:"offline"`
	OutgoingConnectionsCount int    `json:"outgoing_connections_count"`
	RPCConnectionsCount      int    `json:"rpc_connections_count"`
	Stagenet                 bool   `json:"stagenet"`
	StartTime                int    `json:"start_time"`
	Status                   string `json:"status"`
	Synchronized             bool   `json:"synchronized"`
	Target                   int    `json:"target"`
	TargetHeight             int    `json:"target_height"`
	Testnet                  bool   `json:"testnet"`
	TopBlockHash             string `json:"top_block_hash"`
	TxCount                  int    `json:"tx_count"`
	TxPoolSize               int    `json:"tx_pool_size"`
	Untrusted                bool   `json:"untrusted"`
	WasBootstrapEverUsed     bool   `json:"was_bootstrap_ever_used"`
	WhitePeerlistSize        int    `json:"white_peerlist_size"`
}

func (c *Client) GetInfo(ctx context.Context) (*GetInfoResult, error) {
	var resp = &GetInfoResult{}

	if err := c.JsonRPC(ctx, methodGetInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetLastBlockHeaderResult struct {
	BlockHeader struct {
		BlockSize                 int    `json:"block_size"`
		BlockWeight               int    `json:"block_weight"`
		CumulativeDifficulty      int64  `json:"cumulative_difficulty"`
		CumulativeDifficultyTop64 int    `json:"cumulative_difficulty_top64"`
		Depth                     int    `json:"depth"`
		Difficulty                int64  `json:"difficulty"`
		DifficultyTop64           int    `json:"difficulty_top64"`
		Hash                      string `json:"hash"`
		Height                    uint64 `json:"height"`
		LongTermWeight            int    `json:"long_term_weight"`
		MajorVersion              int    `json:"major_version"`
		MinerTxHash               string `json:"miner_tx_hash"`
		MinorVersion              int    `json:"minor_version"`
		Nonce                     int    `json:"nonce"`
		NumTxes                   int    `json:"num_txes"`
		OrphanStatus              bool   `json:"orphan_status"`
		PowHash                   string `json:"pow_hash"`
		PrevHash                  string `json:"prev_hash"`
		Reward                    int64  `json:"reward"`
		Timestamp                 int    `json:"timestamp"`
		WideCumulativeDifficulty  string `json:"wide_cumulative_difficulty"`
		WideDifficulty            string `json:"wide_difficulty"`
	} `json:"block_header"`
}

func (c *Client) GetLastBlockHeader(ctx context.Context) (*GetLastBlockHeaderResult, error) {
	var resp = &GetLastBlockHeaderResult{}

	if err := c.JsonRPC(ctx, methodGetLastBlockHeader, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetCoinbaseTxSumResult struct {
	Credits             int    `json:"credits"`
	EmissionAmount      int64  `json:"emission_amount"`
	EmissionAmountTop64 int    `json:"emission_amount_top64"`
	FeeAmount           int    `json:"fee_amount"`
	FeeAmountTop64      int    `json:"fee_amount_top64"`
	Status              string `json:"status"`
	TopHash             string `json:"top_hash"`
	Untrusted           bool   `json:"untrusted"`
	WideEmissionAmount  string `json:"wide_emission_amount"`
	WideFeeAmount       string `json:"wide_fee_amount"`
}

func (c *Client) GetCoinbaseTxSum(ctx context.Context, height, count uint64) (*GetCoinbaseTxSumResult, error) {
	var (
		resp   = &GetCoinbaseTxSumResult{}
		params = map[string]uint64{
			"height": height,
			"count":  count,
		}
	)

	if err := c.JsonRPC(ctx, methodGetCoinbaseTxSum, params, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetBlockResult struct {
	Blob        string `json:"blob"`
	BlockHeader struct {
		BlockSize                 int    `json:"block_size"`
		BlockWeight               int    `json:"block_weight"`
		CumulativeDifficulty      int64  `json:"cumulative_difficulty"`
		CumulativeDifficultyTop64 int    `json:"cumulative_difficulty_top64"`
		Depth                     int    `json:"depth"`
		Difficulty                int    `json:"difficulty"`
		DifficultyTop64           int    `json:"difficulty_top64"`
		Hash                      string `json:"hash"`
		Height                    int    `json:"height"`
		LongTermWeight            int    `json:"long_term_weight"`
		MajorVersion              int    `json:"major_version"`
		MinerTxHash               string `json:"miner_tx_hash"`
		MinorVersion              int    `json:"minor_version"`
		Nonce                     int    `json:"nonce"`
		NumTxes                   int    `json:"num_txes"`
		OrphanStatus              bool   `json:"orphan_status"`
		PowHash                   string `json:"pow_hash"`
		PrevHash                  string `json:"prev_hash"`
		Reward                    int64  `json:"reward"`
		Timestamp                 int    `json:"timestamp"`
		WideCumulativeDifficulty  string `json:"wide_cumulative_difficulty"`
		WideDifficulty            string `json:"wide_difficulty"`
	} `json:"block_header"`
	Credits     int    `json:"credits"`
	JSON        string `json:"json"`
	MinerTxHash string `json:"miner_tx_hash"`
	Status      string `json:"status"`
	TopHash     string `json:"top_hash"`
	Untrusted   bool   `json:"untrusted"`
}

type GetBlockResultJSON struct {
	MajorVersion int    `json:"major_version"`
	MinorVersion int    `json:"minor_version"`
	Timestamp    int    `json:"timestamp"`
	PrevID       string `json:"prev_id"`
	Nonce        int    `json:"nonce"`
	MinerTx      struct {
		Version    int `json:"version"`
		UnlockTime int `json:"unlock_time"`
		Vin        []struct {
			Gen struct {
				Height int `json:"height"`
			} `json:"gen"`
		} `json:"vin"`
		Vout []struct {
			Amount int64 `json:"amount"`
			Target struct {
				Key string `json:"key"`
			} `json:"target"`
		} `json:"vout"`
		Extra         []int `json:"extra"`
		RctSignatures struct {
			Type int `json:"type"`
		} `json:"rct_signatures"`
	} `json:"miner_tx"`
	TxHashes []string `json:"tx_hashes"`
}

func (j *GetBlockResult) InnerJSON() (*GetBlockResultJSON, error) {
	res := &GetBlockResultJSON{}

	if err := json.Unmarshal([]byte(j.JSON), res); err != nil {
		return nil, nil
	}

	return res, nil
}

func (c *Client) GetBlock(ctx context.Context, height uint64) (*GetBlockResult, error) {
	var (
		resp   = new(GetBlockResult)
		params = map[string]uint64{
			"height": height,
		}
	)

	if err := c.JsonRPC(ctx, methodGetBlock, params, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type GetFeeEstimateResult struct {
	Credits          int    `json:"credits"`
	Fee              int    `json:"fee"`
	QuantizationMask int    `json:"quantization_mask"`
	Status           string `json:"status"`
	TopHash          string `json:"top_hash"`
	Untrusted        bool   `json:"untrusted"`
}

func (c *Client) GetFeeEstimate(ctx context.Context, graceBlocks uint64) (*GetFeeEstimateResult, error) {
	var (
		resp   = new(GetFeeEstimateResult)
		params = map[string]uint64{
			"grace_blocks": graceBlocks,
		}
	)

	if err := c.JsonRPC(ctx, methodGetFeeEstimate, params, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}

type SyncInfoResult struct {
	Credits               int    `json:"credits"`
	Height                int    `json:"height"`
	NextNeededPruningSeed int    `json:"next_needed_pruning_seed"`
	Overview              string `json:"overview"`
	Status                string `json:"status"`
	TargetHeight          int    `json:"target_height"`
	TopHash               string `json:"top_hash"`
	Untrusted             bool   `json:"untrusted"`
	Peers                 []struct {
		Info struct {
			Address           string `json:"address"`
			AddressType       int    `json:"address_type"`
			AvgDownload       int    `json:"avg_download"`
			AvgUpload         int    `json:"avg_upload"`
			ConnectionID      string `json:"connection_id"`
			CurrentDownload   int    `json:"current_download"`
			CurrentUpload     int    `json:"current_upload"`
			Height            int    `json:"height"`
			Host              string `json:"host"`
			Incoming          bool   `json:"incoming"`
			IP                string `json:"ip"`
			LiveTime          int    `json:"live_time"`
			LocalIP           bool   `json:"local_ip"`
			Localhost         bool   `json:"localhost"`
			PeerID            string `json:"peer_id"`
			Port              string `json:"port"`
			PruningSeed       int    `json:"pruning_seed"`
			RecvCount         int    `json:"recv_count"`
			RecvIdleTime      int    `json:"recv_idle_time"`
			RPCCreditsPerHash int    `json:"rpc_credits_per_hash"`
			RPCPort           int    `json:"rpc_port"`
			SendCount         int    `json:"send_count"`
			SendIdleTime      int    `json:"send_idle_time"`
			State             string `json:"state"`
			SupportFlags      int    `json:"support_flags"`
		} `json:"info"`
	} `json:"peers"`
}

func (c *Client) SyncInfo(ctx context.Context) (*SyncInfoResult, error) {
	var resp = new(SyncInfoResult)

	if err := c.JsonRPC(ctx, methodSyncInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}
