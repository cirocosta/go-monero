package main

import (
	"fmt"

	"github.com/cirocosta/go-monero/pkg/rpc"
)

type GetAlternateChainsCommand struct {
	Json bool `long:"json" description:"output the raw json response from the endpoint"`
}

func (c *GetAlternateChainsCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetAlternateChains(ctx)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if c.Json {
		return displayResponseInJSONFormat(resp)
	}

	c.display(resp)
	return nil
}

func (c *GetAlternateChainsCommand) display(v *rpc.GetAlternateChainsResult) {
	table := newTable()

	for _, chain := range v.Chains {
		table.AddRow("Block Hash:", chain.BlockHash)
		table.AddRow("Height:", chain.Height)
		table.AddRow("Main Chain Parent Block:", chain.MainChainParentBlock)
		table.AddRow("Length:", chain.Length)
		table.AddRow("Difficulty:", chain.Difficulty)
		table.AddRow("")
	}

	fmt.Println(table)
}

func init() {
	parser.AddCommand("get-alternate-chains",
		"Get alternate chains",
		"Displays the alternative chains seen by this node",
		&GetAlternateChainsCommand{},
	)
}
