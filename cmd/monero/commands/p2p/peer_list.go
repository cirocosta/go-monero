package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"

	"github.com/cirocosta/go-monero/pkg/levin"
)

type peerListCommand struct {
	NodeAddress string
	Timeout     time.Duration
	Proxy       string
}

func (c *peerListCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peerlist",
		Short: "retrieve a node's peerlist",
		RunE:  c.RunE,
	}

	cmd.Flags().StringVar(&c.NodeAddress, "node-address",
		"", "address of the node to connect to")
	_ = cmd.MarkFlagRequired("node-address")

	cmd.Flags().DurationVar(&c.Timeout, "timeout",
		1*time.Minute, "how long to wait until considering the connection a failure")
	cmd.Flags().StringVar(&c.Proxy, "proxy",
		"", "proxy to proxy connections through (useful for tor)")

	return cmd
}

// nolint:forbidigo
func (c *peerListCommand) RunE(_ *cobra.Command, _ []string) error {
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
	RootCommand.AddCommand((&peerListCommand{}).Cmd())
}
