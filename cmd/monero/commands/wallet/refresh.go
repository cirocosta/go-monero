package wallet

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/wallet"
)

type refreshCommand struct {
	StartHeight uint64

	JSON bool
}

func (c *refreshCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "refresh the wallet openned by the wallet-rpc server",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().Uint64Var(&c.StartHeight, "start-height",
		0, "block height from which to start refreshing")

	return cmd
}

func (c *refreshCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.WalletClient()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.Refresh(ctx, c.StartHeight)
	if err != nil {
		return fmt.Errorf("refresh: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *refreshCommand) pretty(v *wallet.RefreshResult) {
	table := display.NewTable()
	table.AddRow("Blocks Fetched:", v.BlocksFetched)
	table.AddRow("Received Money:", v.ReceivedMoney)
	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&refreshCommand{}).Cmd())
}
