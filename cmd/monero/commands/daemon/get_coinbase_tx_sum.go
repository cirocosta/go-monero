package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getCoinbaseTxSumCommand struct {
	Height uint64
	Count  uint64

	JSON bool
}

func (c *getCoinbaseTxSumCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-coinbase-tx-sum",
		Short: "compute the coinbase amount and the fees amount for n last blocks starting at particular height",
		RunE:  c.RunE,
	}

	cmd.Flags().Uint64Var(&c.Height, "height",
		0, "block height to start the count from")

	cmd.Flags().Uint64Var(&c.Count, "count",
		0, "number of coinbase rewards to include in the sum")

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getCoinbaseTxSumCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetCoinbaseTxSum(ctx, c.Height, c.Count)
	if err != nil {
		return fmt.Errorf("get coinbase tx sum: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *getCoinbaseTxSumCommand) pretty(v *daemon.GetCoinbaseTxSumResult) {
	table := display.NewTable()

	table.AddRow("Emission Amount Top64:", v.EmissionAmountTop64)
	table.AddRow("Emission Amount:", v.EmissionAmount)
	table.AddRow("Fee Amount Top64:", v.FeeAmountTop64)
	table.AddRow("Fee Amount:", v.FeeAmount)
	table.AddRow("Wide Emission Amount:", v.WideEmissionAmount)
	table.AddRow("Wide Fee Amount:", v.WideFeeAmount)

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getCoinbaseTxSumCommand{}).Cmd())
}
