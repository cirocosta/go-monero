package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/daemonrpc"
)

type OnGetBlockHashCommand struct {
	Height uint64 `long:"height" required:"true"`
}

func (c *OnGetBlockHashCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.RequestTimeout)
	defer cancel()

	client, err := daemonrpc.NewClient(options.Address,
		daemonrpc.WithHTTPClient(daemonrpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.OnGetBlockHash(ctx, c.Height)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")

	if err := encoder.Encode(resp); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}

func init() {
	parser.AddCommand("on-get-block-hash",
		"Look up a block's hash by its height",
		"Look up a block's hash by its height",
		&OnGetBlockHashCommand{},
	)
}
