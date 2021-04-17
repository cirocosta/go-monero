package daemonrpc

import "fmt"

const (
	MethodGetBlockCount    = "get_block_count"
	MethodOnGetBlockHash   = "on_get_block_hash"
	MethodGetBlockTemplate = "get_block_template"
	MethodGetConnections   = "get_connections"
	MethodGetInfo          = "get_info"
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

func (c *Client) GetConnections() (*GetConnectionsResult, error) {
	var (
		resp = &GetConnectionsResult{}
	)

	if err := c.JsonRPC(MethodGetConnections, nil, resp); err != nil {
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

func (c *Client) GetInfo() (*GetInfoResult, error) {
	var (
		resp = &GetInfoResult{}
	)

	if err := c.JsonRPC(MethodGetInfo, nil, resp); err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return resp, nil
}
