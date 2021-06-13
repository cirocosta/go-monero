package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/rpc"
)

type GetBlockCommand struct {
	Height *uint64 `long:"height"`
	Hash   *string `long:"hash"`
	Unwrap bool    `long:"unwrap"`

	Json bool `long:"json" description:"output the raw json response from the endpoint"`
}

func (c *GetBlockCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetBlock(ctx, rpc.GetBlockRequestParameters{
		Height: c.Height,
		Hash:   c.Hash,
	})
	if err != nil {
		return fmt.Errorf("get block: %w", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")

	if !c.Unwrap {
		if err := encoder.Encode(resp); err != nil {
			return fmt.Errorf("encode: %w", err)
		}

		return nil
	}

	inner, err := resp.InnerJSON()
	if err != nil {
		return fmt.Errorf("inner json: %w", err)
	}

	if err := encoder.Encode(inner); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}

func init() {
	parser.AddCommand("get-block",
		"Get block",
		"Get block",
		&GetBlockCommand{},
	)
}
