package daemon

// RPCResultFooter contains the set of fields that every RPC result message
// will contain.
//
type RPCResultFooter struct {
	// Status dictates whether the request worked or not. "OK" means good.
	//
	Status string `json:"status"`

	// States if the result is obtained using the bootstrap mode, and is
	// therefore not trusted (`true`), or when the daemon is fully synced
	// and thus handles the RPC locally (`false`).
	//
	Untrusted bool `json:"untrusted"`

	// Credits indicates the number of credits available to the requesting
	// client, if payment for RPC is enabled, otherwise, 0.
	//
	Credits uint64 `json:"credits,omitempty"`

	// TopHash is the hash of the highest block in the chain, If payment
	// for RPC is enabled, otherwise, empty.
	//
	TopHash string `json:"top_hash,omitempty"`
}

// GetAlternateChainsResult is the result of a call to the GetAlternateChains RPC method.
//
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
		Height uint64 `json:"height"`

		// Length is the length in blocks of this alternative chain, after divergence.
		//
		Length uint64 `json:"length"`

		// MainChainParentBlock TODO
		//
		MainChainParentBlock string `json:"main_chain_parent_block"`

		// WideDifficulty is the network difficulty as a hexadecimal
		// string representing a 128-bit number.
		//
		WideDifficulty string `json:"wide_difficulty"`
	} `json:"chains"`

	RPCResultFooter `json:",inline"`
}

// AccessTrackingResult is the result of a call to the RPCAccessTracking RPC method.
//
type RPCAccessTrackingResult struct {
	Data []struct {
		// Count is the number of times that the monero daemon received
		// a request for this RPC method.
		//
		Count uint64 `json:"count"`

		// RPC is the name of the remote procedure call.
		//
		RPC string `json:"rpc"`

		// Time indicates how much time the daemon spent serving this procedure.
		//
		Time uint64 `json:"time"`

		// Credits indicates the number of credits consumed for this method.
		//
		Credits uint64 `json:"credits"`
	} `json:"data"`

	RPCResultFooter `json:",inline"`
}

// HardForkInfoResult is the result of a call to the HardForkInfo RPC method.
//
type HardForkInfoResult struct {
	// EarliestHeight is the earliest height at which <version> is allowed.
	//
	EarliestHeight int `json:"earliest_height"`

	// Whether of not the hard fork is enforced.
	//
	Enabled bool `json:"enabled"`

	// State indicates the current hard fork state:
	//
	// 	0 - likely forked
	//	1 - update needed
	//	2 - ready
	//
	State int `json:"state"`

	// The number of votes required to enable <version>.
	//
	Threshold int `json:"threshold"`

	// Version (<version>) corresponds to the major block version for the
	// fork.
	//
	Version int `json:"version"`

	// Votes is the number of votes to enable <version>
	//
	Votes int `json:"votes"`

	// Voting indicates which version this node is voting for/using.
	//
	Voting int `json:"voting"`

	// Window is the size of the voting window.
	//
	Window int `json:"window"`

	RPCResultFooter `json:",inline"`
}

// GetVersionResult is the result of a call to the GetVersion RPC method.
//
type GetVersionResult struct {
	Release bool   `json:"release"`
	Version uint64 `json:"version"`

	RPCResultFooter `json:",inline"`
}

// GetBansResult is the result of a call to the GetBans RPC method.
//
type GetBansResult struct {
	// Bans contains the list of nodes banned by this node.
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

	RPCResultFooter `json:",inline"`
}

// GetFeeEstimateResult is the result of a call to the GetFeeEstimate RPC
// method.
//
type GetFeeEstimateResult struct {
	// Fee is the per kB fee estimate.
	//
	Fee int `json:"fee"`

	// QuantizationMask indicates that the  fee should be rounded up to an
	// even multiple of this value.
	//
	QuantizationMask int `json:"quantization_mask"`

	RPCResultFooter `json:",inline"`
}

// GetInfoResult is the result of a call to the GetInfo RPC method.
//
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
	Synchronized             bool   `json:"synchronized"`
	Target                   int    `json:"target"`
	TargetHeight             int    `json:"target_height"`
	Testnet                  bool   `json:"testnet"`
	TopBlockHash             string `json:"top_block_hash"`
	TxCount                  int    `json:"tx_count"`
	TxPoolSize               int    `json:"tx_pool_size"`
	WasBootstrapEverUsed     bool   `json:"was_bootstrap_ever_used"`
	WhitePeerlistSize        int    `json:"white_peerlist_size"`

	RPCResultFooter `json:",inline"`
}

// GetBlockTemplateResult is the result of a call to the GetBlockTemplate RPC method.
//
type GetBlockTemplateResult struct {
	// BlockhashingBlob is the blob on which to try to find a valid nonce.
	//
	BlockhashingBlob string `json:"blockhashing_blob"`

	// BlocktemplateBlob is the blob on which to try to mine a new block.
	//
	BlocktemplateBlob string `json:"blocktemplate_blob"`

	// Difficulty is the difficulty of the next block.
	Difficulty int64 `json:"difficulty"`

	// ExpectedReward is the coinbase reward expected to be received if the
	// block is successfully mined.
	//
	ExpectedReward int64 `json:"expected_reward"`

	// Height is the height on which to mine.
	//
	Height int `json:"height"`

	// PrevHash is the hash of the most recent block on which to mine the next block.
	//
	PrevHash string `json:"prev_hash"`

	// ReservedOffset TODO
	//
	ReservedOffset int `json:"reserved_offset"`

	RPCResultFooter `json:",inline"`
}

type Peer struct {
	Host        string `json:"host"`
	ID          uint64 `json:"id"`
	IP          uint32 `json:"ip"`
	LastSeen    int64  `json:"last_seen"`
	Port        uint16 `json:"port"`
	PruningSeed uint32 `json:"pruning_seed"`
	RPCPort     uint16 `json:"rpc_port"`
}

// GetPeerListResult is the result of a call to the GetPeerList RPC method.
//
type GetPeerListResult struct {
	GrayList  []Peer `json:"gray_list"`
	WhiteList []Peer `json:"white_list"`

	RPCResultFooter `json:",inline"`
}

// GetConnectionsResult is the result of a call to the GetConnections RPC method.
//
type GetConnectionsResult struct {
	Connections []struct {
		Address         string `json:"address"`
		AvgDownload     uint64 `json:"avg_download"`
		AvgUpload       uint64 `json:"avg_upload"`
		ConnectionID    string `json:"connection_id"`
		CurrentDownload uint64 `json:"current_download"`
		CurrentUpload   uint64 `json:"current_upload"`
		Height          uint64 `json:"height"`
		Host            string `json:"host"`
		Incoming        bool   `json:"incoming"`
		IP              string `json:"ip"`
		LiveTime        uint64 `json:"live_time"`
		LocalIP         bool   `json:"local_ip"`
		Localhost       bool   `json:"localhost"`
		PeerID          string `json:"peer_id"`
		Port            string `json:"port"`
		RecvCount       uint64 `json:"recv_count"`
		RecvIdleTime    uint64 `json:"recv_idle_time"`
		SendCount       uint64 `json:"send_count"`
		SendIdleTime    uint64 `json:"send_idle_time"`
		State           string `json:"state"`
		SupportFlags    uint64 `json:"support_flags"`
	} `json:"connections"`

	RPCResultFooter `json:",inline"`
}

type GetOutsResult struct {
	Outs []struct {
		Height   uint64 `json:"height"`
		Key      string `json:"key"`
		Mask     string `json:"mask"`
		Txid     string `json:"txid"`
		Unlocked bool   `json:"unlocked"`
	} `json:"outs"`

	RPCResultFooter `json:",inline"`
}

// GetHeightResult is the result of a call to the GetHeight RPC method.
//
type GetHeightResult struct {
	Hash   string `json:"hash"`
	Height uint64 `json:"height"`

	RPCResultFooter `json:",inline"`
}

// GetNetStatsResult is the result of a call to the GetNetStats RPC method.
//
type GetNetStatsResult struct {
	StartTime       int64  `json:"start_time"`
	TotalBytesIn    uint64 `json:"total_bytes_in"`
	TotalBytesOut   uint64 `json:"total_bytes_out"`
	TotalPacketsIn  uint64 `json:"total_packets_in"`
	TotalPacketsOut uint64 `json:"total_packets_out"`

	RPCResultFooter `json:",inline"`
}

// GetPublicNodesResult is the result of a call to the GetPublicNodes RPC method.
//
type GetPublicNodesResult struct {
	WhiteList []Peer `json:"white"`
	GrayList  []Peer `json:"gray"`

	RPCResultFooter `json:",inline"`
}

// GenerateBlocksResult is the result of a call to the GenerateBlocks RPC method.
//
type GenerateBlocksResult struct {
	Blocks []string `json:"blocks"`
	Height int      `json:"height"`

	RPCResultFooter `json:",inline"`
}

// GetBlockCountResult is the result of a call to the GetBlockCount RPC method.
//
type GetBlockCountResult struct {
	Count uint64 `json:"count"`

	RPCResultFooter `json:",inline"`
}

// RelayTxResult is the result of a call to the RelayTx RPC method.
//
type RelayTxResult struct {
	RPCResultFooter `json:",inline"`
}

// GetCoinbaseTxSumResult is the result of a call to the GetCoinbaseTxSum RPC method.
//
type GetCoinbaseTxSumResult struct {
	EmissionAmount      int64  `json:"emission_amount"`
	EmissionAmountTop64 int    `json:"emission_amount_top64"`
	FeeAmount           int    `json:"fee_amount"`
	FeeAmountTop64      int    `json:"fee_amount_top64"`
	WideEmissionAmount  string `json:"wide_emission_amount"`
	WideFeeAmount       string `json:"wide_fee_amount"`

	RPCResultFooter `json:",inline"`
}

type BlockHeader struct {
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
	Height uint64 `json:"height"`

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
	Timestamp int64 `json:"timestamp"`

	// WideCumulativeDifficulty is the cumulative difficulty of all
	// blocks in the blockchain as a hexadecimal string
	// representing a 128-bit number.
	//
	WideCumulativeDifficulty string `json:"wide_cumulative_difficulty"`

	// WideDifficulty is the network difficulty as a hexadecimal
	// string representing a 128-bit number.
	//
	WideDifficulty string `json:"wide_difficulty"`
}

// GetBlockResult is the result of a call to the GetBlock RPC method.
//
type GetBlockResult struct {
	// Blob is a hexadecimal representation of the block.
	//
	Blob string `json:"blob"`

	// BlockHeader contains the details from the block header.
	//
	BlockHeader BlockHeader `json:"block_header"`

	// JSON is a json representation of the block - see `GetBlockResultJSON`.
	//
	JSON string `json:"json"`

	// MinerTxHash is the hash of the coinbase transaction
	//
	MinerTxHash string `json:"miner_tx_hash"`

	RPCResultFooter `json:",inline"`
}

// GetBlockResultJSON is the internal json-formatted block information.
//
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
			Amount uint64 `json:"amount"`
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

func (c *GetBlockResultJSON) MinerOutputs() uint64 {
	res := uint64(0)

	for _, vout := range c.MinerTx.Vout {
		res += vout.Amount
	}

	return res
}

// SyncInfoResult is the result of a call to the SyncInfo RPC method.
//
type SyncInfoResult struct {
	Credits uint64 `json:"credits"`

	Height                uint64 `json:"height"`
	NextNeededPruningSeed uint64 `json:"next_needed_pruning_seed"`
	Overview              string `json:"overview"`
	Status                string `json:"status"`
	TargetHeight          uint64 `json:"target_height"`
	TopHash               string `json:"top_hash"`
	Untrusted             bool   `json:"untrusted"`
	Peers                 []struct {
		Info struct {
			Address           string `json:"address"`
			AddressType       uint64 `json:"address_type"`
			AvgDownload       uint64 `json:"avg_download"`
			AvgUpload         uint64 `json:"avg_upload"`
			ConnectionID      string `json:"connection_id"`
			CurrentDownload   uint64 `json:"current_download"`
			CurrentUpload     uint64 `json:"current_upload"`
			Height            uint64 `json:"height"`
			Host              string `json:"host"`
			IP                string `json:"ip"`
			Incoming          bool   `json:"incoming"`
			LiveTime          uint64 `json:"live_time"`
			LocalIP           bool   `json:"local_ip"`
			Localhost         bool   `json:"localhost"`
			PeerID            string `json:"peer_id"`
			Port              string `json:"port"`
			PruningSeed       uint64 `json:"pruning_seed"`
			RPCCreditsPerHash uint64 `json:"rpc_credits_per_hash"`
			RPCPort           uint64 `json:"rpc_port"`
			RecvCount         uint64 `json:"recv_count"`
			RecvIdleTime      uint64 `json:"recv_idle_time"`
			SendCount         uint64 `json:"send_count"`
			SendIdleTime      uint64 `json:"send_idle_time"`
			State             string `json:"state"`
			SupportFlags      int    `json:"support_flags"`
		} `json:"info"`
	} `json:"peers"`

	RPCResultFooter `json:",inline"`
}

// GetLastBlockHeaderResult is the result of a call to the GetLastBlockHeader RPC method.
//
type GetLastBlockHeaderResult struct {
	BlockHeader BlockHeader `json:"block_header"`

	RPCResultFooter `json:",inline"`
}

// GetBlockHeaderByHeightResult is the result of a call to the GetBlockHeaderByHeight RPC method.
//
type GetBlockHeaderByHeightResult struct {
	BlockHeader BlockHeader `json:"block_header"`

	RPCResultFooter `json:",inline"`
}

// GetBlockHeaderByHashResult is the result of a call to the GetBlockHeaderByHash RPC method.
//
type GetBlockHeaderByHashResult struct {
	BlockHeader  BlockHeader   `json:"block_header"`
	BlockHeaders []BlockHeader `json:"block_headers"`

	RPCResultFooter `json:",inline"`
}

// GetTransactionPoolStatsResult is the result of a call to the GetTransactionPoolStats RPC method.
//
type GetTransactionPoolStatsResult struct {
	PoolStats struct {
		BytesMax   uint64 `json:"bytes_max"`
		BytesMed   uint64 `json:"bytes_med"`
		BytesMin   uint64 `json:"bytes_min"`
		BytesTotal uint64 `json:"bytes_total"`
		FeeTotal   uint64 `json:"fee_total"`
		Histo      []struct {
			Bytes uint64 `json:"bytes"`
			Txs   uint64 `json:"txs"`
		} `json:"histo"`
		Histo98Pc       uint64 `json:"histo_98pc"`
		Num10M          uint64 `json:"num_10m"`
		NumDoubleSpends uint64 `json:"num_double_spends"`
		NumFailing      uint64 `json:"num_failing"`
		NumNotRelayed   uint64 `json:"num_not_relayed"`
		Oldest          int64  `json:"oldest"`
		TxsTotal        uint64 `json:"txs_total"`
	} `json:"pool_stats"`

	RPCResultFooter `json:",inline"`
}

type GetTransactionsResult struct {
	Credits int    `json:"credits"`
	Status  string `json:"status"`
	TopHash string `json:"top_hash"`
	Txs     []struct {
		AsHex           string `json:"as_hex"`
		AsJSON          string `json:"as_json"`
		BlockHeight     uint64 `json:"block_height"`
		BlockTimestamp  int64  `json:"block_timestamp"`
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

type TransactionJSON struct {
	Version    int `json:"version"`
	UnlockTime int `json:"unlock_time"`
	Vin        []struct {
		Key struct {
			Amount     int    `json:"amount"`
			KeyOffsets []uint `json:"key_offsets"`
			KImage     string `json:"k_image"`
		} `json:"key"`
	} `json:"vin"`
	Vout []struct {
		Amount int `json:"amount"`
		Target struct {
			Key string `json:"key"`
		} `json:"target"`
	} `json:"vout"`
	Extra         []byte `json:"extra"`
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

type GetTransactionPoolResult struct {
	Credits        int `json:"credits"`
	SpentKeyImages []struct {
		IDHash    string   `json:"id_hash"`
		TxsHashes []string `json:"txs_hashes"`
	} `json:"spent_key_images"`
	Status       string `json:"status"`
	TopHash      string `json:"top_hash"`
	Transactions []struct {
		BlobSize           uint64 `json:"blob_size"`
		DoNotRelay         bool   `json:"do_not_relay"`
		DoubleSpendSeen    bool   `json:"double_spend_seen"`
		Fee                uint64 `json:"fee"`
		IDHash             string `json:"id_hash"`
		KeptByBlock        bool   `json:"kept_by_block"`
		LastFailedHeight   uint64 `json:"last_failed_height"`
		LastFailedIDHash   string `json:"last_failed_id_hash"`
		LastRelayedTime    uint64 `json:"last_relayed_time"`
		MaxUsedBlockHeight uint64 `json:"max_used_block_height"`
		MaxUsedBlockIDHash string `json:"max_used_block_id_hash"`
		ReceiveTime        int64  `json:"receive_time"`
		Relayed            bool   `json:"relayed"`
		TxBlob             string `json:"tx_blob"`
		TxJSON             string `json:"tx_json"`
		Weight             uint64 `json:"weight"`
	} `json:"transactions"`
	Untrusted bool `json:"untrusted"`
}
