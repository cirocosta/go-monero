package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetFeeEstimateCommand struct {
	GraceBlocks uint64 `long:"grace-blocks" required:"true" description:"TODO"`
}

func (c *GetFeeEstimateCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetFeeEstimate(ctx, c.GraceBlocks)
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
	parser.AddCommand("get-fee-estimate",
		"Gives an estimation on fees per byte",
		"Gives an estimation on fees per byte",
		&GetFeeEstimateCommand{},
	)
}
