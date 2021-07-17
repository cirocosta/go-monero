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

type getBansCommand struct {
	JSON bool
}

func (c *getBansCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-bans",
		Short: "all the nodes that have been banned by our node",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getBansCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetBans(ctx)
	if err != nil {
		return fmt.Errorf("get bans: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

func (c *getBansCommand) pretty(v *daemon.GetBansResult) {
	table := display.NewTable()

	table.AddRow("HOST", "TIME LEFT")

	sort.Slice(v.Bans, func(i, j int) bool {
		return v.Bans[i].Seconds > v.Bans[j].Seconds
	})
	for _, ban := range v.Bans {
		table.AddRow(ban.Host, time.Duration(ban.Seconds)*time.Second)
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getBansCommand{}).Cmd())
}
