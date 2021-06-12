package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GetBlockTemplateCommand struct {
	WalletAddress string `long:"wallet-address" required:"true" description:"address of wallet to receive coinbase txns if block is successfully mined"`
	ReserveSize   uint   `long:"reserve-size" required:"true" description:"reserve size"`
}

func (c *GetBlockTemplateCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
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
