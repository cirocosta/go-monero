package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cirocosta/go-monero/pkg/levin"
)

type P2PHandshakeCommand struct {
	NodeAddress string `long:"node-address" required:"true" description:"address of the node to ping"`
}

func (c *P2PHandshakeCommand) Execute(_ []string) error {
	client, err := levin.NewClient(c.NodeAddress)
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// if err := client.Handshake(ctx); err != nil {
	// 	return fmt.Errorf("ping: %w", err)
	// }

	_, _ = client, ctx

	return nil
}

func init() {
	parser.AddCommand("p2p-handshake",
		"Handshake another node in the p2p network",
		"Handshake another node in the p2p network",
		&P2PHandshakeCommand{},
	)
}
