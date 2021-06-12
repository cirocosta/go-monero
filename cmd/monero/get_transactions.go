package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetTransactionsCommand struct {
	Txns   []string `long:"txn" required:"true"`
	Unwrap bool     `long:"unwrap"`
}

func (c *GetTransactionsCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetTransactions(ctx, c.Txns)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")

	if !c.Unwrap {
		if err := encoder.Encode(resp); err != nil {
			return fmt.Errorf("encode: %w", err)
		}

		return nil
	}

	txns, err := resp.GetTransactions()
	if err != nil {
		return fmt.Errorf("get txns: %w", err)
	}

	if err := encoder.Encode(txns); err != nil {
		return fmt.Errorf("encode txns: %w", err)
	}

	return nil
}

func init() {
	parser.AddCommand("get-transactions",
		"Retrieve transactions",
		"Retrieve transactions",
		&GetTransactionsCommand{},
	)
}
