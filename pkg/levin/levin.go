//
// see https://github.com/monero-project/monero/blob/e45619e61e4831eea70a43fe6985f4d57ea02e9e/contrib/epee/include/net/levin_base.h
// see https://github.com/monero-project/monero/blob/e45619e61e4831eea70a43fe6985f4d57ea02e9e/docs/LEVIN_PROTOCOL.md

package levin

import (
	"encoding/binary"
)

const (
	LevinSignature       uint64 = 0x0101010101012101 // Dander's Nightmare
	LevinProtocolVersion uint32 = 1
	LevinPacketRequest   uint32 = 0x00000001 // Q flag
	LevinPacketReponse   uint32 = 0x00000002 // S flag

	// Return Codes
	LevinOk                               int32 = 0
	LevinErrorConnection                  int32 = -1
	LevinErrorConnectionNotFound          int32 = -2
	LevinErrorConnectionDestroyed         int32 = -3
	LevinErrorConnectionTimedout          int32 = -4
	LevinErrorConnectionNoDuplexProtocol  int32 = -5
	LevinErrorConnectionHandlerNotDefined int32 = -6
	LevinErrorFormat                      int32 = -7
)

const (
	// p2p admin commands
	CommandHandshake    uint32 = 0x1001
	CommandTimedSync    uint32 = 0x1002
	CommandPing         uint32 = 0x1003
	CommandStat         uint32 = 0x1004
	CommandNetworkState uint32 = 0x1005
	CommandPeerID       uint32 = 0x1006
	CommandSupportFlags uint32 = 0x1007
)

//
// Header
//
//
//       0               1               2               3
//       0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7 0 1 2 3 4 5 6 7
//      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//      |      0x01     |      0x21     |      0x01     |      0x01     |
//      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//      |      0x01     |      0x01     |      0x01     |      0x01     |
//      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//      |                             Length                            |
//      |                                                               |
//      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//      |  E. Response  |               _   Command     _
//      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//      		|               _ Return Code   _
//      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//      		|Q|S|B|E|       _       Reserved_
//      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//      		|      0x01     |      0x00     |      0x00     |
//      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//      |     0x00      |
//      +-+-+-+-+-+-+-+-+
//
//
// i.e.,
//
//	BYTE(0X01) BYTE(0X21) BYTE(0X01) BYTE(0X01)  ---.
//							+--> protocol identification
//	BYTE(0X01) BYTE(0X01) BYTE(0X01) BYTE(0X01)  ---'
//
//
//	UINT64(LENGTH)	-----------------------------------> unsigned little-endian 64bit integer
//							     length of the payload _not including_
//							     the header. messages >100MB are rejected.
//
//
//	BYTE(E.RESPONSE) 4BYTE(COMMAND) 4BYTE(RET CODE)
//         |               |		  |
//         |               |		  |
//         |               |	          '->  signed 32-bit little endian integer representing the response
//         |               |		       from the peer from the last command invoked. `0` for request msgs.
//         |               |
//         |               '-> unsigned 32-bit little endian integer
//         |                   representing the monero specific cmd
//         |
//         '-> zero-byte if no response is expected from the peer, non-zero if response is expected.
//	       peers must respond to requests w/ this flag in the same order as received.
//
//
//	BIT(Q) BIT(S) BIT(B) BIT(E) 3BYTE+4BIT(RESERVED)
//         |    |      |      |
//         |    |      |      |
//         |    |      |      '-> set if this is the end of a frag msg
//         |    |      |
//         |    |      '-> set if this is the beginning of a frag msg
//         |    |
//         |    '-> set if the message is a response
//         |
//         '-> set if the message is a request
//
//
//
//	BYTE(0X01) BYTE(0X00) BYTE(0X00) BYTE(0X00)
//         |
//         '--> version
//
type Header struct {
	Signature       uint64
	Length          uint64
	ExpectsResponse bool
	Command         uint32
	ReturnCode      int32
	Flags           uint32 // only 4 most significant bits matter (Q|S|B|E)
	Version         uint32
}

func NewHeader(command uint32, length uint64) *Header {
	return &Header{
		Signature:       LevinSignature,
		Length:          length,
		ExpectsResponse: true,
		Command:         command,
		ReturnCode:      0,
		Flags:           LevinPacketRequest,
		Version:         0x01,
	}
}

func (h *Header) Bytes() []byte {
	var (
		header = make([]byte, 33) // full header
		b      = make([]byte, 8)  // biggest type

		idx  = 0
		size = 0
	)

	{ // signature
		size = 8

		binary.LittleEndian.PutUint64(b, h.Signature)
		copy(header[idx:], b[:size])
		idx += size
	}

	{ // length
		size = 8

		binary.LittleEndian.PutUint64(b, h.Length)
		copy(header[idx:], b[:size])
		idx += size
	}

	{ // expects response
		size = 1

		if h.ExpectsResponse {
			b[0] = 0x01
		} else {
			b[0] = 0x00
		}

		copy(header[idx:], b[:size])
		idx += size
	}

	{ // command
		size = 4

		binary.LittleEndian.PutUint32(b, h.Command)
		copy(header[idx:], b[:size])
		idx += size
	}

	{ // return code
		size = 4

		binary.LittleEndian.PutUint32(b, uint32(h.ReturnCode))
		copy(header[idx:], b[:size])
		idx += size
	}

	{ // return code
		size = 4

		binary.LittleEndian.PutUint32(b, h.Flags)
		copy(header[idx:], b[:size])
		idx += size
	}

	{ // return code
		size = 4

		binary.LittleEndian.PutUint32(b, 0)
		copy(header[idx:], b[:size])
		idx += size
	}

	return header
}
