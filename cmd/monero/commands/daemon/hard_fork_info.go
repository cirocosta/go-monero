package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type hardForkInfoCommand struct {
	JSON bool
}

func (c *hardForkInfoCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hardfork-info",
		Short: "information regarding hard fork voting and readiness.",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *hardForkInfoCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.HardForkInfo(ctx)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

func (c *hardForkInfoCommand) pretty(v *daemon.HardForkInfoResult) {
	table := display.NewTable()

	table.AddRow("Earliest Height:", v.EarliestHeight)
	table.AddRow("Enabled:", v.Enabled)

	state := "unknown"
	switch v.State {
	case 0:
		state = "likely forked"
	case 1:
		state = "update needed"
	case 2:
		state = "ready"
	}

	table.AddRow("State:", state)
	table.AddRow("Threshold:", v.Threshold)
	table.AddRow("Version:", v.Version)
	table.AddRow("Votes:", v.Votes)
	table.AddRow("Voting:", v.Voting)
	table.AddRow("Window:", v.Window)

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&hardForkInfoCommand{}).Cmd())
}
