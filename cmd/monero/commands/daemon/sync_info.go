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

type syncInfoCommand struct {
	JSON bool
}

func (c *syncInfoCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-info",
		Short: "daemon's chain synchronization info",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *syncInfoCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.SyncInfo(ctx)
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
func (c *syncInfoCommand) pretty(v *daemon.SyncInfoResult) {
	table := display.NewTable()

	table.AddRow("HEIGHT", v.Height)
	table.AddRow("")

	sort.Slice(v.Peers, func(i, j int) bool {
		return v.Peers[i].Info.LiveTime > v.Peers[j].Info.LiveTime
	})

	table.AddRow("ADDR", "IN", "STATE", "TIME", "RECV", "SEND")
	for _, peer := range v.Peers {
		table.AddRow(
			peer.Info.Address,
			peer.Info.Incoming,
			peer.Info.State,
			time.Duration(peer.Info.LiveTime)*time.Second,
			humanize.Bytes(peer.Info.RecvCount),
			humanize.Bytes(peer.Info.SendCount),
		)
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&syncInfoCommand{}).Cmd())
}
