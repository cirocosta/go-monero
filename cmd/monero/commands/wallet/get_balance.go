package wallet

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/constant"
	"github.com/cirocosta/go-monero/pkg/rpc/wallet"
)

type getBalanceCommand struct {
	AccountIndex   uint
	AddressIndices []uint
	AllAccounts    bool
	Strict         bool

	JSON bool
}

func (c *getBalanceCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-balance",
		Short: "get the balance of all of the addresses managed by a wallet",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().BoolVar(&c.AllAccounts, "all-accounts",
		false, "retrieve balances from all accounts")

	cmd.Flags().BoolVar(&c.Strict, "strict",
		false, "todo")

	cmd.Flags().UintVar(&c.AccountIndex, "account-index",
		0, "todo")
	cmd.Flags().UintSliceVar(&c.AddressIndices, "address-index",
		[]uint{}, "todo")

	return cmd
}

func (c *getBalanceCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.WalletClient()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	params := wallet.GetBalanceRequestParameters{
		AccountIndex:   c.AccountIndex,
		AddressIndices: c.AddressIndices,
		AllAccounts:    c.AllAccounts,
		Strict:         c.Strict,
	}
	resp, err := client.GetBalance(ctx, params)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

func (c *getBalanceCommand) pretty(v *wallet.GetBalanceResult) {
	c.prettyTotal(v)

	for _, saddr := range v.PerSubaddress {
		c.prettySubAddress(saddr)
	}
}

// nolint:forbidigo
func (c *getBalanceCommand) prettyTotal(v *wallet.GetBalanceResult) {
	table := display.NewTable()

	table.AddRow("Total Balance:", fmt.Sprintf("%f XMR",
		float64(v.Balance)/float64(constant.XMR)))

	if v.BlocksToUnlock > 0 {
		table.AddRow("Total Unlocked Balance:", fmt.Sprintf("%f XMR",
			float64(v.UnlockedBalance)/float64(constant.XMR)))

		table.AddRow("Total Blocks to Unlock:", v.BlocksToUnlock)
		table.AddRow("Total Time to Unlock (s):", v.TimeToUnlock)
	}

	if v.MultisigImportNeeded {
		table.AddRow("Multisig Import Needed:", v.MultisigImportNeeded)
	}

	fmt.Println(table)
}

// nolint:forbidigo
func (c *getBalanceCommand) prettySubAddress(saddr wallet.SubAddress) {
	table := display.NewTable()

	table.AddRow("")
	table.AddRow(saddr.AccountIndex, "Address", saddr.Address)
	table.AddRow("~", "Label:", saddr.Label)
	table.AddRow("~", "UTXOs", saddr.NumUnspentOutputs)
	table.AddRow("~", "Balance:", fmt.Sprintf("%f XMR",
		float64(saddr.Balance)/float64(constant.XMR)))

	if saddr.BlocksToUnlock > 0 {
		table.AddRow("~", "Blocks to Unlock:", saddr.BlocksToUnlock)
		table.AddRow("~", "Time to Unlock:", saddr.TimeToUnlock)
		table.AddRow("~", "Unlocked Balance:", fmt.Sprintf("%f XMR",
			float64(saddr.UnlockedBalance)/float64(constant.XMR)))
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getBalanceCommand{}).Cmd())
}
