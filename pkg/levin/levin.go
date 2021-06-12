//
// see https://github.com/monero-project/monero/blob/e45619e61e4831eea70a43fe6985f4d57ea02e9e/contrib/epee/include/net/levin_base.h
// see https://github.com/monero-project/monero/blob/e45619e61e4831eea70a43fe6985f4d57ea02e9e/docs/LEVIN_PROTOCOL.md

package levin

import (
	"encoding/binary"
	"fmt"
)

const (
	LevinSignature uint64 = 0x0101010101012101 // Dander's Nightmare

	LevinProtocolVersion uint32 = 1

	LevinPacketRequest        uint32 = 0x00000001 // Q flag
	LevinPacketReponse        uint32 = 0x00000002 // S flag
	LevinPacketMaxDefaultSize uint64 = 100000000  // 100MB _after_ handshake
	LevinPacketMaxInitialSize uint64 = 256 * 1024 // 256KiB _before_ handshake

	LevinHeaderSizeBytes = 33
)

const (
	// Return Codes.
	LevinOk                               int32 = 0
	LevinErrorConnection                  int32 = -1
	LevinErrorConnectionNotFound          int32 = -2
	LevinErrorConnectionDestroyed         int32 = -3
	LevinErrorConnectionTimedout          int32 = -4
	LevinErrorConnectionNoDuplexProtocol  int32 = -5
	LevinErrorConnectionHandlerNotDefined int32 = -6
	LevinErrorFormat                      int32 = -7
)

func IsValidReturnCode(c int32) bool {
	// anything >= 0 is good (there are some `1`s in the code :shrug:)
	return c >= LevinErrorFormat
}

const (
	// p2p admin commands.
	CommandHandshake    uint32 = 1001
	CommandTimedSync    uint32 = 1002
	CommandPing         uint32 = 1003
	CommandStat         uint32 = 1004
	CommandNetworkState uint32 = 1005
	CommandPeerID       uint32 = 1006
	CommandSupportFlags uint32 = 1007
)

var (
	MainnetNetworkId = []byte{
		0x12, 0x30, 0xf1, 0x71,
		0x61, 0x04, 0x41, 0x61,
		0x17, 0x31, 0x00, 0x82,
		0x16, 0xa1, 0xa1, 0x10,
	}

	MainnetGenesisTx = "418015bb9ae982a1975da7d79277c2705727a56894ba0fb246adaabb1f4632e3"
)

func IsValidCommand(c uint32) bool {
	return (c >= CommandHandshake && c <= CommandSupportFlags)
}

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

func NewRequestHeader(command uint32, length uint64) *Header {
	return &Header{
		Signature:       LevinSignature,
		Length:          length,
		ExpectsResponse: true,
		Command:         command,
		ReturnCode:      0,
		Flags:           LevinPacketRequest,
		Version:         LevinProtocolVersion,
	}
}

func NewHeaderFromBytesBytes(bytes []byte) (*Header, error) {
	if len(bytes) != LevinHeaderSizeBytes {
		return nil, fmt.Errorf("invalid header size: expected %d, has %d",
			LevinHeaderSizeBytes, len(bytes),
		)
	}

	var (
		size = 0
		idx  = 0
	)

	header := &Header{}

	{ // signature
		size = 8
		header.Signature = binary.LittleEndian.Uint64(bytes[idx : idx+size])
		idx += size

		if header.Signature != LevinSignature {
			return nil, fmt.Errorf("signature mismatch: expected %x, got %x",
				LevinSignature, header.Signature,
			)
		}
	}

	{ // length
		size = 8
		header.Length = binary.LittleEndian.Uint64(bytes[idx : idx+size])
		idx += size
	}

	{ // expects response
		size = 1
		header.ExpectsResponse = (bytes[idx] != 0)
		idx += size
	}

	{ // command
		size = 4
		header.Command = binary.LittleEndian.Uint32(bytes[idx : idx+size])
		idx += size

		if !IsValidCommand(header.Command) {
			return nil, fmt.Errorf("invalid command %d", header.Command)
		}
	}

	{ // return code
		size = 4
		header.ReturnCode = int32(binary.LittleEndian.Uint32(bytes[idx : idx+size]))
		idx += size

		if !IsValidReturnCode(header.ReturnCode) {
			return nil, fmt.Errorf("invalid return code %d", header.ReturnCode)
		}
	}

	{ // flags
		size = 4
		header.Flags = binary.LittleEndian.Uint32(bytes[idx : idx+size])
		idx += size
	}

	{ // version
		size = 4
		header.Version = binary.LittleEndian.Uint32(bytes[idx : idx+size])
		idx += size

		if header.Version != LevinProtocolVersion {
			return nil, fmt.Errorf("invalid version %x",
				header.Version)
		}
	}

	return header, nil
}

func (h *Header) Bytes() []byte {
	var (
		header = make([]byte, LevinHeaderSizeBytes) // full header
		b      = make([]byte, 8)                    // biggest type

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

	{ // flags
		size = 4

		binary.LittleEndian.PutUint32(b, h.Flags)
		copy(header[idx:], b[:size])
		idx += size
	}

	{ // version
		size = 4

		binary.LittleEndian.PutUint32(b, h.Version)
		copy(header[idx:], b[:size])
		idx += size
	}

	return header
}
