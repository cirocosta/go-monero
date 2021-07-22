package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getFeeEstimateCommand struct {
	GraceBlocks uint64
	JSON        bool
}

func (c *getFeeEstimateCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-fee-estimate",
		Short: "estimate fees in atomic units per kB",
		RunE:  c.RunE,
	}

	cmd.Flags().Uint64Var(&c.GraceBlocks, "grace-blocks",
		10, "number of blocks we want the fee to be valid for")

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getFeeEstimateCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetFeeEstimate(ctx, c.GraceBlocks)
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
func (c *getFeeEstimateCommand) pretty(v *daemon.GetFeeEstimateResult) {
	table := display.NewTable()

	table.AddRow("Fee:", v.Fee)
	table.AddRow("Quantization Mask:", v.QuantizationMask)

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getFeeEstimateCommand{}).Cmd())
}
