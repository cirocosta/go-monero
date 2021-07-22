package wallet

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/wallet"
)

type createAddressCommand struct {
	AccountIndex uint
	Count        uint
	Label        string

	JSON bool
}

func (c *createAddressCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-address",
		Short: "create a new address for an account, optionally labelling it",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().UintVar(&c.AccountIndex, "account-index",
		0, "account to create the address for")
	cmd.Flags().StringVar(&c.Label, "label",
		"", "label for the new address")
	cmd.Flags().UintVar(&c.Count, "count",
		1, "number of addresses to create")

	return cmd
}

func (c *createAddressCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.WalletClient()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.CreateAddress(ctx, c.AccountIndex, c.Count, c.Label)
	if err != nil {
		return fmt.Errorf("create addr: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *createAddressCommand) pretty(v *wallet.CreateAddressResult) {
	table := display.NewTable()

	for idx := range v.AddressIndices {
		table.AddRow(idx, "Address:", v.Addresses[idx])
		table.AddRow("~", "Address Index:", v.AddressIndices[idx])

		if idx != len(v.AddressIndices)-1 {
			table.AddRow("")
		}
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&createAddressCommand{}).Cmd())
}
