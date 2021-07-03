package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GenerateBlocksCommand struct {
	AmountOfBlocks uint64 `long:"amount-of-blocks" default:"1"`
	WalletAddress  string `long:"wallet-address" required:"true"`

	PreviousBlock string `long:"previous-block"`
	StartingNonce uint32 `long:"starting-nonce"`
}

func (c *GenerateBlocksCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	// TODO make use of arguments
	//
	resp, err := client.GenerateBlocks(ctx)
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
		&GenerateBlocksCommand{},
	)
}
