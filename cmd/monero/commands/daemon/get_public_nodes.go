package daemon

import (
	"fmt"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getPublicNodesCommand struct {
	Gray           bool
	White          bool
	IncludeBlocked bool

	JSON bool
}

func (c *getPublicNodesCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-public-nodes",
		Short: "all known peers advertising as public nodes",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.Gray, "gray",
		false, "whether or not to include peers from gray peerlist")

	cmd.Flags().BoolVar(&c.White, "white",
		true, "whether or not to include peers from white peerlist")

	cmd.Flags().BoolVar(&c.IncludeBlocked, "include-blocked",
		false, "whether or not to include blocked")

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getPublicNodesCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetPublicNodes(ctx, daemon.GetPublicNodesRequestParameters{
		Gray:           c.Gray,
		White:          c.White,
		IncludeBlocked: c.IncludeBlocked,
	})
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *getPublicNodesCommand) pretty(v *daemon.GetPublicNodesResult) {
	table := display.NewTable()

	table.AddRow("TYPE", "HOST", "RPC PORT", "SINCE")
	for _, peer := range v.GrayList {
		table.AddRow("GRAY", peer.Host, peer.RPCPort, "")
	}

	sort.Slice(v.WhiteList, func(i, j int) bool {
		return v.WhiteList[i].LastSeen < v.WhiteList[j].LastSeen
	})
	for _, peer := range v.WhiteList {
		table.AddRow("WHITE", peer.Host, peer.RPCPort, humanize.Time(time.Unix(peer.LastSeen, 0)))
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getPublicNodesCommand{}).Cmd())
}
