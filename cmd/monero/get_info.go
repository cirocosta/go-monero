package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/daemonrpc"
)

type GetInfoCommand struct{}

func (c *GetInfoCommand) Execute(_ []string) error {
	client, err := daemonrpc.NewClient(options.Address,
		daemonrpc.WithHTTPClient(daemonrpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.GetInfo()
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
