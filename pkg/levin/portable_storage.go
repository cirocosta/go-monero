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
	Serializable Serializable `json:"-,omitempty"`
	Value        interface{}
}

func (e Entry) String() string {
	v, ok := e.Value.(string)
	if !ok {
		panic(fmt.Errorf("interface couldnt be casted to string"))
	}

	return v
}

func (e Entry) Uint8() uint8 {
	v, ok := e.Value.(uint8)
	if !ok {
		panic(fmt.Errorf("interface couldnt be casted to uint8"))
	}

	return v
}

func (e Entry) Uint16() uint16 {
	v, ok := e.Value.(uint16)
	if !ok {
		panic(fmt.Errorf("interface couldnt be casted to uint16"))
	}

	return v
}

func (e Entry) Uint32() uint32 {
	v, ok := e.Value.(uint32)
	if !ok {
		panic(fmt.Errorf("interface couldnt be casted to uint32"))
	}

	return v
}

func (e Entry) Uint64() uint64 {
	v, ok := e.Value.(uint64)
	if !ok {
		panic(fmt.Errorf("interface couldnt be casted to uint64"))
	}

	return v
}

func (e Entry) Entries() Entries {
	v, ok := e.Value.(Entries)
	if !ok {
		panic(fmt.Errorf("interface couldnt be casted to levin.Entries"))
	}

	return v
}

func (e Entry) Bytes() []byte {
	return nil
}

type Entries []Entry

func (e Entries) Bytes() []byte {
	return nil
}

type PortableStorage struct {
	Entries Entries
}

func NewPortableStorageFromBytes(bytes []byte) (*PortableStorage, error) {
	var (
		size = 0
		idx  = 0
	)

	{ // sig-a
		size = 4

		if len(bytes[idx:]) < size {
			return nil, fmt.Errorf("sig-a out of bounds")
		}

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

	ps := &PortableStorage{}

	_, ps.Entries = ReadObject(bytes[idx:])

	return ps, nil
}

func ReadString(bytes []byte) (int, string) {
	idx := 0

	n, strLen := ReadVarInt(bytes)
	idx += n

	return idx + strLen, string(bytes[idx : idx+strLen])
}

func ReadObject(bytes []byte) (int, Entries) {
	idx := 0

	n, i := ReadVarInt(bytes[idx:])
	idx += n

	entries := make(Entries, i)

	for iter := 0; iter < i; iter++ {
		entries[iter] = Entry{}
		entry := &entries[iter]

		lenName := int(bytes[idx])
		idx += 1

		entry.Name = string(bytes[idx : idx+lenName])
		idx += lenName

		ttype := bytes[idx]
		idx += 1

		n, obj := ReadAny(bytes[idx:], ttype)
		idx += n

		entry.Value = obj
	}

	return idx, entries
}

func ReadArray(ttype byte, bytes []byte) (int, Entries) {
	var (
		idx = 0
		n   = 0
	)

	n, i := ReadVarInt(bytes[idx:])
	idx += n

	entries := make(Entries, i)

	for iter := 0; iter < i; iter++ {
		n, obj := ReadAny(bytes[idx:], ttype)
		idx += n

		entries[iter] = Entry{
			Value: obj,
		}
	}

	return idx, entries
}

func ReadAny(bytes []byte, ttype byte) (int, interface{}) {
	var (
		idx = 0
		n   = 0
	)

	if ttype&BoostSerializeFlagArray != 0 {
		internalType := ttype &^ BoostSerializeFlagArray
		n, obj := ReadArray(internalType, bytes[idx:])
		idx += n

		return idx, obj
	}

	if ttype == BoostSerializeTypeObject {
		n, obj := ReadObject(bytes[idx:])
		idx += n

		return idx, obj
	}

	if ttype == BoostSerializeTypeUint8 {
		obj := uint8(bytes[idx])
		n += 1
		idx += n

		return idx, obj
	}

	if ttype == BoostSerializeTypeUint16 {
		obj := binary.LittleEndian.Uint16(bytes[idx:])
		n += 2
		idx += n

		return idx, obj
	}

	if ttype == BoostSerializeTypeUint32 {
		obj := binary.LittleEndian.Uint32(bytes[idx:])
		n += 4
		idx += n

		return idx, obj
	}

	if ttype == BoostSerializeTypeUint64 {
		obj := binary.LittleEndian.Uint64(bytes[idx:])
		n += 8
		idx += n

		return idx, obj
	}

	if ttype == BoostSerializeTypeInt64 {
		obj := binary.LittleEndian.Uint64(bytes[idx:])
		n += 8
		idx += n

		return idx, int64(obj)
	}

	if ttype == BoostSerializeTypeString {
		n, obj := ReadString(bytes[idx:])
		idx += n

		return idx, obj
	}

	panic(fmt.Errorf("unknown ttype %x", ttype))
	return -1, nil
}

// reads var int, returning number of bytes read and the integer in that byte
// sequence.
//
func ReadVarInt(b []byte) (int, int) {
	sizeMask := b[0] & PortableRawSizeMarkMask

	switch uint32(sizeMask) {
	case uint32(PortableRawSizeMarkByte):
		return 1, int(b[0] >> 2)
	case uint32(PortableRawSizeMarkWord):
		return 2, int((binary.LittleEndian.Uint16(b[0:2])) >> 2)
	case PortableRawSizeMarkDword:
		return 4, int((binary.LittleEndian.Uint32(b[0:4])) >> 2)
	case uint32(PortableRawSizeMarkInt64):
		panic("int64 not supported") // TODO
		// return int((binary.LittleEndian.Uint64(b[0:8])) >> 2)
		//         '-> bad
	default:
		panic(fmt.Errorf("malformed sizemask: %+v", sizeMask))
	}

	return -1, -1
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
