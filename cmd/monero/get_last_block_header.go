package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetLastBlockHeaderCommand struct {
}

func (c *GetLastBlockHeaderCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetLastBlockHeader(ctx)
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
	parser.AddCommand("get-last-block-header",
		"Get the header of the last block",
		"Get the header of the last block",
		&GetLastBlockHeaderCommand{},
	)
}
