package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cirocosta/go-monero/pkg/levin"
)

var (
	fpath = flag.String("f", "resp.bin", "file location")
)

func run() error {
	flag.Parse()

	f, err := os.Open(*fpath)
	if err != nil {
		return fmt.Errorf("open %s: %w", *fpath, err)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	ps, err := levin.NewPortableStorageFromBytes(b)
	if err != nil {
		return fmt.Errorf("portable storage from bytes: %w", err)
	}

	pl := levin.NewLocalPeerListFromEntries(ps.Entries)

	for addr := range pl.Peers {
		fmt.Println(addr)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
