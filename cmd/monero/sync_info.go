package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type SyncInfoCommand struct{}

func (c *SyncInfoCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.SyncInfo(ctx)
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
	parser.AddCommand("sync-info",
		"Get synchronisation information (restricted)",
		"Get synchronisation information (restricted)",
		&SyncInfoCommand{},
	)
}
