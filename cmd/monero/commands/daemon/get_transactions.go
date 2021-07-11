package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
)

type getTransactionsCommand struct {
	Txns   []string
	Unwrap bool
}

func (c *getTransactionsCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-transactions",
		Short: "lookup one or more transactions by hash",
		RunE:  c.RunE,
	}

	cmd.Flags().StringArrayVar(&c.Txns, "txn",
		[]string{}, "hash of a transaction to lookup")
	cmd.MarkFlagRequired("txn")

	cmd.Flags().BoolVar(&c.Unwrap, "unwrap",
		false, "whether or not to unwrap the json representation of the transaction")

	return cmd
}

func (c *getTransactionsCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetTransactions(ctx, c.Txns)
	if err != nil {
		return fmt.Errorf("get transactions: %w", err)
	}

	if !c.Unwrap {
		return display.JSON(resp)
	}

	txns, err := resp.GetTransactions()
	if err != nil {
		return fmt.Errorf("resp get txns: %w", err)
	}

	return display.JSON(txns)
}

func init() {
	RootCommand.AddCommand((&getTransactionsCommand{}).Cmd())
}
