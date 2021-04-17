package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cirocosta/go-monero/pkg/levin"
)

type P2PPingCommand struct {
	NodeAddress string `long:"node-address" required:"true" description:"address of the node to ping"`
}

func (c *P2PPingCommand) Execute(_ []string) error {
	client, err := levin.NewClient(c.NodeAddress)
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	return nil
}

func init() {
	parser.AddCommand("p2p-ping",
		"Ping another node in the p2p network",
		"Ping another node in the p2p network",
		&P2PPingCommand{},
	)
}
