package main

import (
	"encoding/binary"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/jessevdk/go-flags"
)

var opts = struct {
	Uint16 uint16 `long:"uint16"`
	Uint32 uint32 `long:"uint32"`
	String string `long:"string"`
}{}

func main() {
	parser := flags.NewParser(&opts, flags.Default)

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	var b []byte

	switch {
	case opts.Uint16 != 0:
		b = make([]byte, 2)
		binary.LittleEndian.PutUint16(b, opts.Uint16)
	case opts.Uint32 != 0:
		b = make([]byte, 4)
		binary.LittleEndian.PutUint32(b, opts.Uint32)
	case len(opts.String) != 0:
		b = []byte(opts.String)
	default:
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	spew.Dump(b)
}
