package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cirocosta/go-monero/pkg/levin"
)

var (
	fpath = flag.String("f", "", "file location")
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

	fmt.Println(ps)
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
