package levin

import (
	"encoding/binary"
	"fmt"
)

const (
	PortableStorageSignatureA    uint32 = 0x01011101
	PortableStorageSignatureB    uint32 = 0x01020101
	PortableStorageFormatVersion byte   = 0x01

	PortableRawSizeMarkMask  byte   = 0x03
	PortableRawSizeMarkByte  byte   = 0x00
	PortableRawSizeMarkWord  uint16 = 0x01
	PortableRawSizeMarkDword uint32 = 0x02
	PortableRawSizeMarkInt64 uint64 = 0x03
)

type Entry struct {
	Name         string
	Serializable Serializable
}

//
// ex: node_data:
//       foo: bar
//     payload_data
//       caz: 1
//
type PortableStorage struct {
	Entries []Entry
}

func NewPortableStorageFromBytes(bytes []byte) (*PortableStorage, error) {

	var (
		size = 0
		idx  = 0
	)

	{ // sig-a
		size = 4
		sig := binary.LittleEndian.Uint32(bytes[idx : idx+size])
		idx += size

		if sig != uint32(PortableStorageSignatureA) {
			return nil, fmt.Errorf("sig-a doesn't match")
		}
	}

	{ // sig-b
		size = 4
		sig := binary.LittleEndian.Uint32(bytes[idx : idx+size])
		idx += size

		if sig != uint32(PortableStorageSignatureB) {
			return nil, fmt.Errorf("sig-b doesn't match")
		}
	}

	{ // format ver
		size = 1
		version := bytes[idx]
		idx += size

		if version != PortableStorageFormatVersion {
			return nil, fmt.Errorf("version doesn't match")
		}
	}

	// read section
	//	count = read var int (gives you number of entries)
	//	while count >0:
	//

	return nil, nil
}

func ReadVarInt(b []byte) int {
	sizeMask := b[0] & PortableRawSizeMarkMask

	switch uint32(sizeMask) {
	case uint32(PortableRawSizeMarkByte):
		return int(b[0] >> 2)
	case uint32(PortableRawSizeMarkWord):
		return int((binary.LittleEndian.Uint16(b[0:2])) >> 2)
	case PortableRawSizeMarkDword:
		return int((binary.LittleEndian.Uint32(b[0:4])) >> 2)
	case uint32(PortableRawSizeMarkInt64):
		panic("int64 not supported") // TODO
		// return int((binary.LittleEndian.Uint64(b[0:8])) >> 2)
		//         '-> bad
	default:
		panic(fmt.Errorf("malformed sizemask: %+v", sizeMask))
	}

	return 0
}

func (s *PortableStorage) Bytes() []byte {
	var (
		body = make([]byte, 9) // fit _at least_ signatures + format ver
		b    = make([]byte, 8) // biggest type

		idx  = 0
		size = 0
	)

	{ // signature a
		size = 4

		binary.LittleEndian.PutUint32(b, PortableStorageSignatureA)
		copy(body[idx:], b[:size])
		idx += size
	}

	{ // signature b
		size = 4

		binary.LittleEndian.PutUint32(b, PortableStorageSignatureB)
		copy(body[idx:], b[:size])
		idx += size
	}

	{ // format ver
		size = 1

		b[0] = PortableStorageFormatVersion
		copy(body[idx:], b[:size])
		idx += size
	}

	// // write_var_in
	varInB, err := VarIn(len(s.Entries))
	if err != nil {
		panic(fmt.Errorf("varin '%d': %w", len(s.Entries), err))
	}

	body = append(body, varInB...)
	for _, entry := range s.Entries {
		fmt.Println("xx", entry.Name)
		body = append(body, byte(len(entry.Name))) // section name length
		body = append(body, []byte(entry.Name)...) // section name
		body = append(body, entry.Serializable.Bytes()...)
	}

	return body
}

type Serializable interface {
	Bytes() []byte
}

type Section struct {
	Entries []Entry
}

func (s Section) Bytes() []byte {
	body := []byte{
		BoostSerializeTypeObject,
	}

	varInB, err := VarIn(len(s.Entries))
	if err != nil {
		panic(fmt.Errorf("varin '%d': %w", len(s.Entries), err))
	}

	body = append(body, varInB...)
	for _, entry := range s.Entries {
		fmt.Println("  xx", entry.Name)
		body = append(body, byte(len(entry.Name))) // section name length
		body = append(body, []byte(entry.Name)...) // section name
		body = append(body, entry.Serializable.Bytes()...)
	}

	return body
}

func VarIn(i int) ([]byte, error) {
	if i <= 63 {
		return []byte{
			(byte(i) << 2) | PortableRawSizeMarkByte,
		}, nil
	}

	if i <= 16383 {
		b := []byte{0x00, 0x00}
		binary.LittleEndian.PutUint16(b,
			(uint16(i)<<2)|PortableRawSizeMarkWord,
		)

		return b, nil
	}

	if i <= 1073741823 {
		b := []byte{0x00, 0x00, 0x00, 0x00}
		binary.LittleEndian.PutUint32(b,
			(uint32(i)<<2)|PortableRawSizeMarkDword,
		)

		return b, nil
	}

	return nil, fmt.Errorf("int %d too big", i)
}
