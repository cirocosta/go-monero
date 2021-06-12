package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/rpc"
)

type RelayTxCommand struct {
	Txns []string `long:"txn" required:"true" description:"transaction to relay"`
}

func (c *RelayTxCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.RequestTimeout)
	defer cancel()

	client, err := rpc.NewClient(options.Address,
		rpc.WithHTTPClient(rpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.RelayTx(ctx, c.Txns)
	if err != nil {
		return fmt.Errorf("relay tx: %w", err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "  ")

	if err := encoder.Encode(resp); err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	return nil
}

func init() {
	parser.AddCommand("relay-tx",
		"Relay txns",
		"Relay txns",
		&RelayTxCommand{},
	)
}
