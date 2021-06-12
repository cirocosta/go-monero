package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type RelayTxCommand struct {
	Txns []string `long:"txn" required:"true" description:"transaction to relay"`
}

func (c *RelayTxCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
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
