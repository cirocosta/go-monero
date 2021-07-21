package wallet

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/wallet"
)

type getHeightCommand struct {
	JSON bool
}

func (c *getHeightCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-height",
		Short: "get the wallet's current block height",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getHeightCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.WalletClient()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetHeight(ctx)
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
func (c *getHeightCommand) pretty(v *wallet.GetHeightResult) {
	table := display.NewTable()

	table.AddRow("Height:", v.Height)

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getHeightCommand{}).Cmd())
}
