package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type stopMiningCommand struct {
	MinerAddress     string
	BackgroundMining bool
	IgnoreBattery    bool
	ThreadsCount     uint

	JSON bool
}

func (c *stopMiningCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-mining",
		Short: "stop mining on the daemon",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *stopMiningCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.StopMining(ctx)
	if err != nil {
		return fmt.Errorf("stop mining: %w", err)
	}
	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *stopMiningCommand) pretty(v *daemon.StopMiningResult) {
	table := display.NewTable()
	table.AddRow("Status:", v.Status)
	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&stopMiningCommand{}).Cmd())
}
