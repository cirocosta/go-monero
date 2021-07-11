package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
)

type getTransactionPoolCommand struct{}

func (c *getTransactionPoolCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-transaction-pool",
		Short: "information about valid transactions seen by the node but not yet mined into a block, including spent key image info for the txpool",
		RunE:  c.RunE,
	}

	return cmd
}

func (c *getTransactionPoolCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetTransactionPool(ctx)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	return display.JSON(resp)
}

func init() {
	RootCommand.AddCommand((&getTransactionPoolCommand{}).Cmd())
}
