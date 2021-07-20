package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type miningStatusCommand struct {
	JSON bool
}

func (c *miningStatusCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mining-status",
		Short: "information about this daemon's mining activity",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *miningStatusCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.MiningStatus(ctx)
	if err != nil {
		return fmt.Errorf("get bans: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *miningStatusCommand) pretty(v *daemon.MiningStatusResult) {
	table := display.NewTable()

	table.AddRow("Active:", v.Active)
	table.AddRow("Address:", v.Address)
	table.AddRow("Background Idle Threshold:", v.BgIdleThreshold)
	table.AddRow("Background Ignore Battery:", v.BgIgnoreBattery)
	table.AddRow("Background Minimum Idle Seconds:", v.BgMinIdleSeconds)
	table.AddRow("Background Target:", v.BgTarget)
	table.AddRow("Block Reward:", v.BlockReward)
	table.AddRow("Block Target:", v.BlockTarget)
	table.AddRow("Difficulty:", v.Difficulty)
	table.AddRow("Difficulty Top64:", v.DifficultyTop64)
	table.AddRow("Background Mining Enabled:", v.IsBackgroundMiningEnabled)
	table.AddRow("Proof-of-Work Algorithm:", v.PowAlgorithm)
	table.AddRow("Speed", v.Speed)
	table.AddRow("Threads:", v.ThreadsCount)
	table.AddRow("Wide Difficulty:", v.WideDifficulty)

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&miningStatusCommand{}).Cmd())
}
