package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/rpc"
)

type GetCoinbaseTxSumCommand struct {
	Height uint64 `long:"height" required:"true" description:"block height from which getting the amounts"`
	Count  uint64 `long:"count" required:"true" description:"number of blocks to include in the sum"`
}

func (c *GetCoinbaseTxSumCommand) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.RequestTimeout)
	defer cancel()

	client, err := rpc.NewClient(options.Address,
		rpc.WithHTTPClient(rpc.NewHTTPClient(options.Verbose)),
	)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", options.Address, err)
	}

	resp, err := client.GetCoinbaseTxSum(ctx, c.Height, c.Count)
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
	parser.AddCommand("get-coinbase-tx-sum",
		"Get the coinbase amount and the fees amount for n last blocks starting at particular height",
		"Get the coinbase amount and the fees amount for n last blocks starting at particular height",
		&GetCoinbaseTxSumCommand{},
	)
}
