package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cirocosta/go-monero/pkg/rpc"
)

type GetPublicNodesCommand struct {
	Gray           bool `long:"gray"`
	White          bool `long:"white"`
	IncludeBlocked bool `long:"include-blocked"`
}

func (c *GetPublicNodesCommand) Execute(_ []string) error {
	ctx, cancel := options.Context()
	defer cancel()

	client, err := options.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetPublicNodes(ctx, rpc.GetPublicNodesRequestParameters{
		Gray:           c.Gray,
		White:          c.White,
		IncludeBlocked: c.IncludeBlocked,
	})
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
	parser.AddCommand("get-public-nodes",
		"Get public nodes",
		"Retrieves the list of public nodes known by the node - gray/white/banned will only be included if set by flags",
		&GetPublicNodesCommand{},
	)
}
