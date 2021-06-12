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
	methodRPCAccessTracking  = "rpc_access_tracking"
	methodRelayTx            = "relay_tx"
	methodSyncInfo           = "sync_info"
)

type GetAlternateChainsResult struct {
	// Chains is the array of alternate chains seen by the node.
	//
	Chains []struct {
		// BlockHash is the hash of the first diverging block of this alternative chain.
		//
		BlockHash string `json:"block_hash"`

		// BlockHashes TODO
		//
		BlockHashes []string `json:"block_hashes"`

		// Difficulty is the cumulative difficulty of all blocks in the alternative chain.
		//
		Difficulty int64 `json:"difficulty"`

		// DifficultyTop64 is the most-significat 64 bits of the
		// 128-bit network difficulty.
		//
		DifficultyTop64 int `json:"difficulty_top64"`

		// Height is the block height of the first diverging block of this alternative chain.
		//
		Height int `json:"height"`

		// Length is the length in blocks of this alternative chain, after divergence.
		//
		Length int `json:"length"`

		// MainChainParentBlock TODO
		//
		MainChainParentBlock string `json:"main_chain_parent_block"`

		// WideDifficulty is the network difficulty as a hexadecimal
		// string representing a 128-bit number.
		//
		WideDifficulty string `json:"wide_difficulty"`
	} `json:"chains"`

	// Status dictates whether the request worked or not. "OK" means good.
	Status string `json:"status"`

	// States if the result is obtained using the bootstrap mode, and is
	// therefore not trusted (`true`), or when the daemon is fully synced
	// and thus handles the RPC locally (`false`).
	Untrusted bool `json:"untrusted"`
}

// GetAlternateChains displays alternative chains seen by the node.
//
// (restricted)
//
func (c *Client) GetAlternateChains(ctx context.Context) (*GetAlternateChainsResult, error) {
	resp := &GetAlternateChainsResult{}

	if err := c.JSONRPC(ctx, methodGetAlternateChains, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

type AccessTrackingResult struct {
	Data []struct {
		Count   uint64 `json:"count"`
		Credits uint64 `json:"credits"`
		RPC     string `json:"rpc"`
		Time    uint64 `json:"time"`
	} `json:"data"`
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
}

func (c *Client) RPCAccessTracking(ctx context.Context) (*AccessTrackingResult, error) {
	resp := &AccessTrackingResult{}

	if err := c.JSONRPC(ctx, methodRPCAccessTracking, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
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
	resp := &HardForkInfoResult{}

	if err := c.JSONRPC(ctx, methodHardForkInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

type GetBansResult struct {
	// Bans contains the list of banned nodes.
	//
	Bans []struct {
		// Host is the string representation of the node that is banned.
		//
		Host string `json:"host"`

		// IP is the integer representation of the host banned.
		//
		IP int `json:"ip"`

		// Seconds represents how many seconds are left for the ban to
		// be lifted.
		//
		Seconds uint `json:"seconds"`
	} `json:"bans"`
	Status    string `json:"status"`
	Untrusted bool   `json:"untrusted"`
}

// GetBans retrieves the list of banned IPs.
//
// (restrited)
//
func (c *Client) GetBans(ctx context.Context) (*GetBansResult, error) {
	resp := &GetBansResult{}

	if err := c.JSONRPC(ctx, methodGetBans, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

type GetBlockCountResult struct {
	Count  uint64 `json:"count"`
	Status string `json:"status"`
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

	if err := c.JSONRPC(ctx, methodRelayTx, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
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

	if err := c.JSONRPC(ctx, methodGetBlockTemplate, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
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
	resp := &GetConnectionsResult{}

	if err := c.JSONRPC(ctx, methodGetConnections, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
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
	resp := &GetInfoResult{}

	if err := c.JSONRPC(ctx, methodGetInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
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
	resp := &GetLastBlockHeaderResult{}

	if err := c.JSONRPC(ctx, methodGetLastBlockHeader, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
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

	if err := c.JSONRPC(ctx, methodGetCoinbaseTxSum, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}

type GetBlockResult struct {
	// Blob is a hexadecimal representation of the block.
	//
	Blob        string `json:"blob"`
	BlockHeader struct {
		// BlockSize is the block size in bytes.
		//
		BlockSize uint64 `json:"block_size"`

		// BlockWeight TODO
		//
		BlockWeight uint64 `json:"block_weight"`

		// CumulativeDifficulty is the cumulative difficulty of all
		// blocks up to this one.
		//
		CumulativeDifficulty uint64 `json:"cumulative_difficulty"`
		// CumulativeDifficultyTop64 most significant 64 bits of the
		// 128-bit cumulative difficulty.
		//
		CumulativeDifficultyTop64 uint64 `json:"cumulative_difficulty_top64"`
		// Depth is the number of blocks succeeding this block on the
		// blockchain. (the larger this number, the oldest this block
		// is).
		//
		Depth uint64 `json:"depth"`
		// Difficulty is the difficulty that was set for mining this block.
		//
		Difficulty uint64 `json:"difficulty"`
		// DifficultyTop64 corresponds to the most significat 64-bit of
		// the 128-bit difficulty.
		//
		DifficultyTop64 uint64 `json:"difficulty_top64"`
		// Hash is the hash of this block.
		//
		Hash string `json:"hash"`
		// Height is the number of blocks preceding this block on the blockchain.
		//
		Height uint `json:"height"`
		// LongTermWeight TODO
		//
		LongTermWeight uint64 `json:"long_term_weight"`
		// MajorVersion is the major version of the monero protocol at
		// this block height.
		//
		MajorVersion uint `json:"major_version"`
		// MinerTxHash TODO
		//
		MinerTxHash string `json:"miner_tx_hash"`
		// MinorVersion is the minor version of the monero protocol at
		// this block height.
		//
		MinorVersion uint `json:"minor_version"`
		// Nonce is the cryptographic random one-time number used in
		// mining this block.
		//
		Nonce uint64 `json:"nonce"`
		// NumTxes is the number of transactions in this block, not
		// counting the coinbase tx.
		//
		NumTxes uint `json:"num_txes"`
		// OrphanStatus indicates whether this block is part of the
		// longest chain or not (true == not part of it).
		//
		OrphanStatus bool `json:"orphan_status"`
		// PowHash TODO
		//
		PowHash string `json:"pow_hash"`
		// PrevHash is the hash of the block immediately preceding this
		// block in the chain.
		//
		PrevHash string `json:"prev_hash"`
		// Reward the amount of new atomic-units generated in this
		// block and rewarded to the miner (1XMR = 1e12 atomic units).
		//
		Reward uint64 `json:"reward"`
		// Timestamp is the unix timestamp at which the block was
		// recorded into the blockchain.
		//
		Timestamp uint64 `json:"timestamp"`
		// WideCumulativeDifficulty is the cumulative difficulty of all
		// blocks in the blockchain as a hexadecimal string
		// representing a 128-bit number.
		//
		WideCumulativeDifficulty string `json:"wide_cumulative_difficulty"`
		// WideDifficulty is the network difficulty as a hexadecimal
		// string representing a 128-bit number.
		//
		WideDifficulty string `json:"wide_difficulty"`
	} `json:"block_header"`
	Credits int `json:"credits"`
	// JSON is a json representation of the block - see `GetBlockResultJSON`.
	//
	JSON        string `json:"json"`
	MinerTxHash string `json:"miner_tx_hash"`
	Status      string `json:"status"`
	TopHash     string `json:"top_hash"`
	Untrusted   bool   `json:"untrusted"`
}

type GetBlockResultJSON struct {
	// MajorVersion (same as in the block header)
	//
	MajorVersion uint `json:"major_version"`

	// MinorVersion (same as in the block header)
	//
	MinorVersion uint `json:"minor_version"`

	// Timestamp (same as in the block header)
	//
	Timestamp uint64 `json:"timestamp"`

	// PrevID (same as `block_hash` in the block header)
	//
	PrevID string `json:"prev_id"`

	// Nonce (same as in the block header)
	//
	Nonce int `json:"nonce"`

	// MinerTx contains the miner transaction information.
	//
	MinerTx struct {
		// Version is the transaction version number
		//
		Version int `json:"version"`

		// UnlockTime is the block height when the coinbase transaction becomes spendable.
		//
		UnlockTime int `json:"unlock_time"`

		// Vin lists the transaction inputs.
		//
		Vin []struct {
			Gen struct {
				Height int `json:"height"`
			} `json:"gen"`
		} `json:"vin"`

		// Vout lists the transaction outputs.
		//
		Vout []struct {
			Amount int64 `json:"amount"`
			Target struct {
				Key string `json:"key"`
			} `json:"target"`
		} `json:"vout"`
		// Extra (aka the transaction id) can be used to include any
		// random 32byte/64char hex string.
		//
		Extra []int `json:"extra"`

		// RctSignatures contain the signatures of tx signers.
		//
		// ps.: coinbase txs DO NOT have signatures.
		//
		RctSignatures struct {
			Type int `json:"type"`
		} `json:"rct_signatures"`
	} `json:"miner_tx"`

	// TxHashes is the list of hashes of non-coinbase transactions in the
	// block.
	//
	TxHashes []string `json:"tx_hashes"`
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

// GetBlockParameters represents the set of possible parameters that can be used
// for submitting a call to the `get_block` jsonrpc method.
//
type GetBlockParameters struct {
	Height *uint64 `json:"height"`
	Hash   *string `json:"string"`
}

func (p GetBlockParameters) Validate() error {
	if p.Height == nil && p.Hash == nil {
		return fmt.Errorf("height or hash must be set")
	}

	if p.Height != nil && p.Hash != nil {
		return fmt.Errorf("either height or hash must be set, not both")
	}

	return nil
}

// GetBlock fetches full block information from a block at a particular hash OR height.
//
func (c *Client) GetBlock(ctx context.Context, params GetBlockParameters) (*GetBlockResult, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	resp := &GetBlockResult{}
	if err := c.JSONRPC(ctx, methodGetBlock, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
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

	if err := c.JSONRPC(ctx, methodGetFeeEstimate, params, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
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
	resp := new(SyncInfoResult)

	if err := c.JSONRPC(ctx, methodSyncInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("jsonrpc: %w", err)
	}

	return resp, nil
}
