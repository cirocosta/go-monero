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

type getPeerListCommand struct {
	JSON bool

	White bool
	Gray  bool
}

func (c *getPeerListCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-peer-list",
		Short: "peers lists (white and gray)",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().BoolVar(&c.White, "white",
		true, "whether or not to show the white list")
	cmd.Flags().BoolVar(&c.Gray, "gray",
		false, "whether or not show gray list")

	return cmd
}

func (c *getPeerListCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	if !c.White && !c.Gray {
		return fmt.Errorf("either white or gray (or both) must be set")
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

// nolint:forbidigo
func (c *getPeerListCommand) pretty(v *daemon.GetPeerListResult) {
	table := display.NewTable()

	table.AddRow("TYPE", "HOST", "PORT ", "RPC", "SINCE")

	sort.Slice(v.WhiteList, func(i, j int) bool {
		return v.WhiteList[i].LastSeen < v.WhiteList[j].LastSeen
	})

	if c.Gray {
		for _, peer := range v.GrayList {
			table.AddRow("gray",
				peer.Host, peer.Port,
				peer.RPCPort, "")
		}
	}

	if c.White {
		for _, peer := range v.WhiteList {
			table.AddRow("white",
				peer.Host, peer.Port,
				peer.RPCPort,
				humanize.Time(time.Unix(peer.LastSeen, 0)),
			)
		}

	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getPeerListCommand{}).Cmd())
}
