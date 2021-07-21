package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type startMiningCommand struct {
	MinerAddress     string
	BackgroundMining bool
	IgnoreBattery    bool
	ThreadsCount     uint

	JSON bool
}

func (c *startMiningCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-mining",
		Short: "start mining on the daemon",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().StringVar(&c.MinerAddress, "address",
		"", "address to send coinbase transactions to")
	_ = cmd.MarkFlagRequired("address")

	cmd.Flags().UintVar(&c.ThreadsCount, "threads",
		1, "number of threads to dedicate to mining")

	cmd.Flags().BoolVar(&c.BackgroundMining, "background-mining",
		false, "whether the miner should run in the background or "+
			"foreground")

	cmd.Flags().BoolVar(&c.IgnoreBattery, "ignore-battery",
		true, "if laptop battery state should be ignore or not")

	return cmd
}

func (c *startMiningCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	params := daemon.StartMiningRequestParameters{
		MinerAddress:     c.MinerAddress,
		BackgroundMining: c.BackgroundMining,
		IgnoreBattery:    c.IgnoreBattery,
		ThreadsCount:     c.ThreadsCount,
	}
	resp, err := client.StartMining(ctx, params)
	if err != nil {
		return fmt.Errorf("start mining: %w", err)
	}
	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *startMiningCommand) pretty(v *daemon.StartMiningResult) {
	table := display.NewTable()
	table.AddRow("Status:", v.Status)
	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&startMiningCommand{}).Cmd())
}
