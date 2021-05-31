package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/daemonrpc"
)

type GetBlockCommand struct {
	Height uint64 `long:"height" required:"true"`
	Unwrap bool   `long:"unwrap"`
}

func (c *GetBlockCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.RequestTimeout)
	defer cancel()

	client, err := daemonrpc.NewClient(options.Address,
		daemonrpc.WithHTTPClient(daemonrpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.GetBlock(ctx, c.Height)
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
