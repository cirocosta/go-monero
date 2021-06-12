package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetBlockCountCommand struct{}

func (c *GetBlockCountCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetBlockCount(ctx)
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
	parser.AddCommand("get-block-count",
		"Get the block count",
		"Look up how many blocks are in the longest chain known to the node",
		&GetBlockCountCommand{},
	)
}
