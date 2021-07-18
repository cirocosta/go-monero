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

type rpcAccessTrackingCommand struct {
	JSON bool
}

func (c *rpcAccessTrackingCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rpc-access-tracking",
		Short: "statistics about rpc access",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *rpcAccessTrackingCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.RPCAccessTracking(ctx)
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
func (c *rpcAccessTrackingCommand) pretty(v *daemon.RPCAccessTrackingResult) {
	table := display.NewTable()

	sort.Slice(v.Data, func(i, j int) bool {
		return v.Data[i].Time < v.Data[j].Time
	})

	table.AddRow("RPC", "COUNT", "TIME SPENT SERVING", "CREDITS")
	for _, entry := range v.Data {
		table.AddRow(entry.RPC, entry.Count, time.Duration(entry.Time), entry.Credits)
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&rpcAccessTrackingCommand{}).Cmd())
}
