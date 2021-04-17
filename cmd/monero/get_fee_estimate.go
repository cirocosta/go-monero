package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/daemonrpc"
)

type GetFeeEstimateCommand struct {
	GraceBlocks uint64 `long:"grace-blocks" required:"true" description:"TODO"`
}

func (c *GetFeeEstimateCommand) Execute(_ []string) error {
	client, err := daemonrpc.NewClient(options.Address,
		daemonrpc.WithHTTPClient(daemonrpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.GetFeeEstimate(c.GraceBlocks)
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
