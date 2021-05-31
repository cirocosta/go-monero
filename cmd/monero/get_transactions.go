package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/daemonrpc"
)

type GetTransactionsCommand struct {
	Txns   []string `long:"txn" required:"true"`
	Unwrap bool     `long:"unwrap" default:"true"`
}

func (c *GetTransactionsCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.RequestTimeout)
	defer cancel()

	client, err := daemonrpc.NewClient(options.Address,
		daemonrpc.WithHTTPClient(daemonrpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.GetTransactions(ctx, c.Txns)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if !c.Unwrap {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent(" ", "  ")

		if err := encoder.Encode(resp); err != nil {
			return fmt.Errorf("encode: %w", err)
		}

		return nil
	}

	txns, err := resp.GetTransactions()
	if err != nil {
		return fmt.Errorf("get txns: %w", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")

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
