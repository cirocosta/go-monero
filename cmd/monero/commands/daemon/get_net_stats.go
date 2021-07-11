package daemon

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getNetStatsCommand struct {
	JSON bool
}

func (c *getNetStatsCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-net-stats",
		Short: "networking statistics.",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getNetStatsCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetNetStats(ctx)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

func (c *getNetStatsCommand) pretty(v *daemon.GetNetStatsResult) {
	table := display.NewTable()

	table.AddRow("Start Time", time.Unix(v.StartTime, 0))
	table.AddRow("Total In", display.ByteSize(v.TotalBytesIn))
	table.AddRow("Total Out", display.ByteSize(v.TotalBytesOut))
	table.AddRow("Total Packets In", v.TotalPacketsIn)
	table.AddRow("Total Packets Out", v.TotalPacketsOut)

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getNetStatsCommand{}).Cmd())
}
