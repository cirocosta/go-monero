package daemon

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getPeerListCommand struct {
	JSON bool
}

func (c *getPeerListCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-peer-list",
		Short: "peers lists (white and gray)",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getPeerListCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetPeerList(ctx)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

func (c *getPeerListCommand) pretty(v *daemon.GetPeerListResult) {
	table := display.NewTable()

	table.AddRow("TYPE", "HOST", "PORT ", "RPC", "LAST SEEN")

	for _, peer := range v.GrayList {
		table.AddRow("GRAY", peer.Host, peer.Port, peer.RPCPort, "")
	}

	sort.Slice(v.WhiteList, func(i, j int) bool {
		return v.WhiteList[i].LastSeen < v.WhiteList[j].LastSeen
	})
	for _, peer := range v.WhiteList {
		table.AddRow("WHITE", peer.Host, peer.Port, peer.RPCPort, display.Since(time.Unix(peer.LastSeen, 0)))
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getPeerListCommand{}).Cmd())
}
