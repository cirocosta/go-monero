package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cirocosta/go-monero/pkg/levin"
)

type P2PPeerList struct {
	NodeAddress string `long:"node-address" default:"xps.utxo.com.br:18080" description:"address of the node to find the peer list of"`
}

func (c *P2PPeerList) Execute(_ []string) error {
	client, err := levin.NewClient(c.NodeAddress)
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	pl, err := client.Handshake(ctx)
	if err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	for addr := range pl.Peers {
		fmt.Println(addr)
	}

	return nil
}

func init() {
	parser.AddCommand("p2p-peer-list",
		"Find out the list of local peers known by a node",
		"Find out the list of local peers known by a node",
		&P2PPeerList{},
	)
}
