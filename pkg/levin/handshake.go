package levin

// Go representation of COMMAND_HANDSHAKE_T
//
// see p2p_protocol_defs.h
//
type CmdHandshake struct {
	NodeData    BasicNodeData
	PayloadData PayloadData
}

// Go representation of basic_node_data
//
// see p2p_protocol_defs.h
//
type BasicNodeData struct {
	NetworkId         []byte
	MyPort            uint32
	RPCPort           uint16
	RPCCreditsPerHash uint32
	PeerId            uint64
	SupportFlags      uint32
}

func (d *BasicNodeData) Bytes() []byte {
	return nil
}

type PayloadData struct {
	CumulativeDifficulty uint64
	CurrentHeight        uint64
	TopId                string
	TopVersion           byte
}

func (d *PayloadData) Bytes() []byte {
	return nil
}
