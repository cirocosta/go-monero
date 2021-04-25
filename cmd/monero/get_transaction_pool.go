package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/daemonrpc"
)

type GetTransactionPoolCommand struct{}

func (c *GetTransactionPoolCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.RequestTimeout)
	defer cancel()

	client, err := daemonrpc.NewClient(options.Address,
		daemonrpc.WithHTTPClient(daemonrpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.GetTransactionPool(ctx)
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
	parser.AddCommand("get-transaction-pool",
		"Get all transactions in the pool",
		"Show information about valid transactions seen by the node but not yet mined into a block, as well as spent key image information for the txpool in the node's memory.",
		&GetTransactionPoolCommand{},
	)
}
