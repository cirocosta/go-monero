package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetTransactionPoolCommand struct{}

func (c *GetTransactionPoolCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
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
