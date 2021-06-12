package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/rpc"
)

type GetInfoCommand struct{}

func (c *GetInfoCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.RequestTimeout)
	defer cancel()

	client, err := rpc.NewClient(options.Address,
		rpc.WithHTTPClient(rpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.GetInfo(ctx)
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
	parser.AddCommand("get-info",
		"Retrieve general information about the state of your node and the network. (restricted)",
		"Retrieve general information about the state of your node and the network. (restricted)",
		&GetInfoCommand{},
	)
}
