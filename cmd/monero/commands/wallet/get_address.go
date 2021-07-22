package wallet

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/wallet"
)

type getAddressCommand struct {
	AccountIndex   uint
	AddressIndices []uint

	JSON bool
}

func (c *getAddressCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-address",
		Short: "addresses for an account",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().UintVar(&c.AccountIndex, "account-index",
		0, "todo")
	cmd.Flags().UintSliceVar(&c.AddressIndices, "address-index",
		[]uint{}, "todo")

	return cmd
}

func (c *getAddressCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.WalletClient()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	params := wallet.GetAddressRequestParameters{
		AccountIndex:   c.AccountIndex,
		AddressIndices: c.AddressIndices,
	}
	resp, err := client.GetAddress(ctx, params)
	if err != nil {
		return fmt.Errorf("get address: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *getAddressCommand) pretty(v *wallet.GetAddressResult) {
	table := display.NewTable()

	table.AddRow("", "ADDRESS", "LABEL", "USED")
	for _, addr := range v.Addresses {
		table.AddRow(
			addr.AddressIndex, addr.Address, addr.Label, addr.Used,
		)
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getAddressCommand{}).Cmd())
}
