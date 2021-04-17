package levin

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/lunixbochs/struc"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Signature [8]byte

var MoneroSignature = Signature{
	0x01, 0x21, 0x01, 0x01,
	0x01, 0x01, 0x01, 0x01,
}

type Command uint32

const (
	// p2p admin commands
	CommandHandshake    Command = 0x1001
	CommandTimedSync    Command = 0x1002
	CommandPing         Command = 0x1003
	CommandStat         Command = 0x1004
	CommandNetworkState Command = 0x1005
	CommandPeerID       Command = 0x1006
	CommandSupportFlags Command = 0x1007
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
	Signature        Signature
	Length           uint64 `struc:"uint64,little"`
	ExpectedResponse uint8
	Command          Command `struc:"uint32,little"`
	ReturnCode       uint32  `struc:"uint32,little"`
	Flags            byte    // only 4 most significant bits matter (Q|S|B|E)
	Reserved         [3]byte `struc:"pad"`
	Version          byte
	Padding          [3]byte `struc:"pad"`
}

func (h *Header) Bytes() []byte {
	var buf bytes.Buffer

	if err := struc.Pack(&buf, h); err != nil {
		panic(fmt.Errorf("pack: %w", err)) // this should never fail
	}

	return buf.Bytes()
}

type Request struct {
	Header  Header
	Payload []byte
}

var CommandPayloads = map[Command][]byte{
	CommandHandshake:    []byte("handshake"),
	CommandNetworkState: []byte("network_state"),
	CommandPeerID:       []byte("peer_id"),
	CommandPing:         []byte("ping"),
	CommandStat:         []byte("stat"),
	CommandSupportFlags: []byte("support_flags"),
	CommandTimedSync:    []byte("timed_sync"),
}

func NewRequest(command Command) *Request {
	payload, found := CommandPayloads[command]
	if !found {
		panic("programming mistake: map doesn't have command '%s'")
	}

	header := Header{
		Signature:        MoneroSignature,
		Length:           uint64(len(payload)),
		ExpectedResponse: uint8(rand.Uint32()),
		Command:          command,
		ReturnCode:       0,
		Flags:            1 << 7, // Q
		Reserved:         [3]byte{0, 0, 0},
		Version:          0x01,
		Padding:          [3]byte{0, 0, 0},
	}

	return &Request{
		Header:  header,
		Payload: payload,
	}
}

func (r *Request) Bytes() []byte {
	return append(r.Header.Bytes(), r.Payload...)
}
