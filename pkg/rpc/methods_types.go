package rpc

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
	Credits uint64 `json:"credits"`
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

	// TopHash is the hash of the highest block in the chain, If payment
	// for RPC is enabled, otherwise, empty.
	//
	TopHash string `json:"top_hash,omitempty"`

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

// TODO
type GenerateBlocksResult struct {
}
