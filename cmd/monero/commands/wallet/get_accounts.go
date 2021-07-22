package wallet

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/wallet"
)

type getAccountsCommand struct {
	Tag            string
	StrictBalances bool

	JSON bool
}

func (c *getAccountsCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-accounts",
		Short: "retrieve wallet's accounts",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().StringVar(&c.Tag, "tag",
		"", "only display accounts with this tag")
	cmd.Flags().BoolVar(&c.StrictBalances, "strict-balances",
		false, "todo")

	return cmd
}

func (c *getAccountsCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.WalletClient()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	params := wallet.GetAccountsRequestParameters{
		Tag:            c.Tag,
		StrictBalances: c.StrictBalances,
	}
	resp, err := client.GetAccounts(ctx, params)
	if err != nil {
		return fmt.Errorf("get accounts: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

func (c *getAccountsCommand) pretty(v *wallet.GetAccountsResult) {
	c.prettySummary(v)
	c.prettyAccounts(v)
}

// nolint:forbidigo
func (c *getAccountsCommand) prettySummary(v *wallet.GetAccountsResult) {
	table := display.NewTable()

	table.AddRow("Total Balance:", display.XMR(v.TotalBalance))
	table.AddRow("Total Unlocked Balance:", display.XMR(v.TotalUnlockedBalance))
	table.AddRow("")

	fmt.Println(table)
}

// nolint:forbidigo
func (c *getAccountsCommand) prettyAccounts(v *wallet.GetAccountsResult) {
	table := display.NewTable()

	table.AddRow("", "LABEL", "TAG", "ADDR", "BALANCE", "UNLOCKED BALANCE")
	for _, account := range v.SubaddressAccounts {
		table.AddRow(
			account.AccountIndex,
			account.Label,
			account.Tag,
			options.RootOpts.AddrFmter()(account.BaseAddress),
			display.XMR(account.Balance),
			display.XMR(account.UnlockedBalance),
		)
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getAccountsCommand{}).Cmd())
}
