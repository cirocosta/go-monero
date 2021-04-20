package main

import (
	"encoding/binary"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
)

var opts = struct {
	Byte   byte   `long:"byte"`
	Uint16 uint16 `long:"uint16"`
	Uint32 uint32 `long:"uint32"`
	String string `long:"string"`

	HexToDecimal string `long:"hex-to-decimal"`
}{}

func main() {
	parser := flags.NewParser(&opts, flags.Default)

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	var b []byte

	switch {
	case opts.Byte != 0:
		b = []byte{opts.Byte}
	case opts.Uint16 != 0:
		b = make([]byte, 2)
		binary.LittleEndian.PutUint16(b, opts.Uint16)
	case opts.Uint32 != 0:
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, opts.Uint32)
	case len(opts.String) != 0:
		b = []byte(opts.String)
	case len(opts.HexToDecimal) != 0:
		i, err := strconv.ParseUint(opts.HexToDecimal, 16, 64)
		if err != nil {
			panic(err)
		}

		b = make([]byte, 8)
		spew.Dump(i)
		binary.LittleEndian.PutUint64(b, i)
	default:
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	spew.Dump(b)
}
