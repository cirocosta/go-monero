package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type OnGetBlockHashCommand struct {
	Height uint64 `long:"height" required:"true"`
}

func (c *OnGetBlockHashCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.OnGetBlockHash(ctx, c.Height)
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
	parser.AddCommand("on-get-block-hash",
		"Look up a block's hash by its height",
		"Look up a block's hash by its height",
		&OnGetBlockHashCommand{},
	)
}
