package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getOutsCommand struct {
	Outputs []uint
	GetTxID bool
	JSON    bool
}

func (c *getOutsCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-outs",
		Short: "output details",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().BoolVar(&c.GetTxID, "get-txid",
		true, "include the transaction id in the response")

	cmd.Flags().UintSliceVar(&c.Outputs, "output",
		[]uint{}, "key offsets to lookup output information about")
	cmd.MarkFlagRequired("output")

	return cmd
}

func (c *getOutsCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetOuts(ctx, c.Outputs, c.GetTxID)
	if err != nil {
		return fmt.Errorf("get outs: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

func (c *getOutsCommand) pretty(v *daemon.GetOutsResult) {
	table := display.NewTable()

	table.AddRow("HEIGHT", "KEY", "TXID", "UNLOCKED")
	for _, out := range v.Outs {
		table.AddRow(out.Height, out.Key, out.Txid, out.Unlocked)
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getOutsCommand{}).Cmd())
}
