package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetNetStatsCommand struct{}

func (c *GetNetStatsCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetNetStats(ctx)
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
	parser.AddCommand("get-net-stats",
		"Get network statistics",
		"Get network statistics",
		&GetNetStatsCommand{},
	)
}
