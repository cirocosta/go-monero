package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/rpc"
)

type GetBlockHeaderByHashCommand struct {
	Hashes []string `long:"hashes"`
	Hash   *string  `long:"hash"`
	Unwrap bool     `long:"unwrap"`

	Json bool `long:"json" description:"output the raw json response from the endpoint"`
}

func (c *GetBlockHeaderByHashCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetBlockHeaderByHash(ctx, rpc.GetBlockHeaderByHashRequestParameters{
		Hashes: c.Hashes,
		Hash:   c.Hash,
	})
	if err != nil {
		return fmt.Errorf("get block: %w", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")

	if err := encoder.Encode(resp); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}

func init() {
	parser.AddCommand("get-block-header-by-hash",
		"Get block header by hash",
		"Get block header by hash",
		&GetBlockHeaderByHashCommand{},
	)
}
