package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetTransactionPoolStatsCommand struct {
}

func (c *GetTransactionPoolStatsCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetTransactionPoolStats(ctx)
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
	parser.AddCommand("get-transaction-pool-stats",
		"Get the transaction pool statistics",
		"Get the transaction pool statistics",
		&GetTransactionPoolStatsCommand{},
	)
}
