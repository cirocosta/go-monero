package zmq

type Topic string

const (
	TopicUnknown Topic = "unknown"

	TopicMinimalTxPoolAdd Topic = "json-minimal-txpool_add"
	TopicFullTxPoolAdd    Topic = "json-full-txpool_add"

	TopicMinimalChainMain Topic = "json-minimal-chain_main"
	TopicFullChainMain    Topic = "json-full-chain_main"
)

type MinimalChainMain struct {
	FirstHeight uint64   `json:"first_height"`
	FirstPrevID string   `json:"first_prev_id"`
	Ids         []string `json:"ids"`
}

type FullChainMain struct {
	MajorVersion int    `json:"major_version"`
	MinorVersion int    `json:"minor_version"`
	Timestamp    int64  `json:"timestamp"`
	PrevID       string `json:"prev_id"`
	Nonce        uint64 `json:"nonce"`
	MinerTx      struct {
		Version    int   `json:"version"`
		UnlockTime int64 `json:"unlock_time"`
		Inputs     []struct {
			Gen struct {
				Height uint64 `json:"height"`
			} `json:"gen"`
		} `json:"inputs"`
		Outputs []struct {
			Amount uint64 `json:"amount"`
			ToKey  struct {
				Key string `json:"key"`
			} `json:"to_key"`
		} `json:"outputs"`
		Extra      string        `json:"extra"`
		Signatures []interface{} `json:"signatures"`
		Ringct     struct {
			Type        int           `json:"type"`
			Encrypted   []interface{} `json:"encrypted"`
			Commitments []interface{} `json:"commitments"`
			Fee         uint64        `json:"fee"`
		} `json:"ringct"`
	} `json:"miner_tx"`
	TxHashes []string `json:"tx_hashes"`
}

type MinimalTxPoolAdd struct {
	ID       string `json:"id"`
	BlobSize uint64 `json:"blob_size"`
}

type FullTxPoolAdd struct {
	Version    int   `json:"version"`
	UnlockTime int64 `json:"unlock_time"`
	Inputs     []struct {
		ToKey struct {
			Amount     uint64   `json:"amount"`
			KeyOffsets []uint64 `json:"key_offsets"`
			KeyImage   string   `json:"key_image"`
		} `json:"to_key"`
	} `json:"inputs"`
	Outputs []struct {
		Amount int `json:"amount"`
		ToKey  struct {
			Key string `json:"key"`
		} `json:"to_key"`
	} `json:"outputs"`
	Extra      string        `json:"extra"`
	Signatures []interface{} `json:"signatures"`
	Ringct     struct {
		Type      int `json:"type"`
		Encrypted []struct {
			Mask   string `json:"mask"`
			Amount string `json:"amount"`
		} `json:"encrypted"`
		Commitments []string `json:"commitments"`
		Fee         int      `json:"fee"`
		Prunable    struct {
			RangeProofs  []interface{} `json:"range_proofs"`
			Bulletproofs []struct {
				V      []string `json:"V"`
				AUpper string   `json:"A"`
				S      string   `json:"S"`
				T1     string   `json:"T1"`
				T2     string   `json:"T2"`
				Taux   string   `json:"taux"`
				Mu     string   `json:"mu"`
				L      []string `json:"L"`
				R      []string `json:"R"`
				ALower string   `json:"a"`
				B      string   `json:"b"`
				T      string   `json:"t"`
			} `json:"bulletproofs"`
			Mlsags     []interface{} `json:"mlsags"`
			PseudoOuts []string      `json:"pseudo_outs"`
		} `json:"prunable"`
	} `json:"ringct"`
}
