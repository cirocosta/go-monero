package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cirocosta/go-monero/pkg/levin"
	"golang.org/x/net/proxy"
)

type P2PPeerList struct {
	NodeAddress string        `long:"node-address" default:"xps.utxo.com.br:18080" description:"address of the node to find the peer list of"`
	Timeout     time.Duration `long:"timeout" default:"20s"`
	Proxy       string        `long:"proxy" short:"x" description:"socks5 proxy addr"`
}

func (c *P2PPeerList) Execute(_ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	opts := []levin.ClientOption{}

	if c.Proxy != "" {
		dialer, err := proxy.SOCKS5("tcp", c.Proxy, nil, nil)
		if err != nil {
			return fmt.Errorf("socks5 '%s': %w", c.Proxy, err)
		}

		contextDialer, ok := dialer.(proxy.ContextDialer)
		if !ok {
			panic("can't cast proxy dialer to proxy context dialer")
		}

		opts = append(opts, levin.WithContextDialer(contextDialer))
	}

	client, err := levin.NewClient(ctx, c.NodeAddress, opts...)
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	defer client.Close()

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
