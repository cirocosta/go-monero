package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetInfoCommand struct{}

func (c *GetInfoCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
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
