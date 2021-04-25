package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/daemonrpc"
)

type GetBlockTemplateCommand struct {
	WalletAddress string `long:"wallet-address" required:"true" description:"address of wallet to receive coinbase txns if block is successfully mined"`
	ReserveSize   uint   `long:"reserve-size" required:"true" description:"reserve size"`
}

func (c *GetBlockTemplateCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.RequestTimeout)
	defer cancel()

	client, err := daemonrpc.NewClient(options.Address,
		daemonrpc.WithHTTPClient(daemonrpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.GetBlockTemplate(ctx, c.WalletAddress, c.ReserveSize)
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
	parser.AddCommand("get-block-template",
		"Get a block template on which mining a new block",
		"Get a block template on which mining a new block",
		&GetBlockTemplateCommand{},
	)
}
