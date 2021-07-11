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

	JSON bool
}

func (c *getBalanceCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-balance",
		Short: "get the balance of all of the addresses managed by a wallet",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(
		&c.JSON,
		"json",
		false,
		"whether or not to output the result as json",
	)

	return cmd
}

func (c *getBalanceCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.WalletClient()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetBalance(ctx)
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
	table := display.NewTable()

	table.AddRow("BALANCE", fmt.Sprintf("%f XMR",
		float64(v.Balance)/float64(constant.XMR)))

	if v.BlocksToUnlock > 0 {
		table.AddRow("UNLOCKED BALANCE", fmt.Sprintf("%f XMR",
			float64(v.UnlockedBalance)/float64(constant.XMR)))

		table.AddRow("BLOCKS TO UNLOCK", v.BlocksToUnlock)
		table.AddRow("TIME TO UNLOCK (s)", v.TimeToUnlock)
	}

	if v.MultisigImportNeeded {
		table.AddRow("MULTISIG IMPORT NEEDED", v.MultisigImportNeeded)
	}

	table.AddRow("")

	for _, saddr := range v.PerSubaddress {
		table.AddRow("", "ACCOUNT IDX", saddr.AccountIndex)
		table.AddRow("", "ADDRESS IDX", saddr.AddressIndex)
		table.AddRow("", "ADDRESS", saddr.Address)
		table.AddRow("", "LABEL", saddr.Label)
		table.AddRow("", "UTXOs", saddr.NumUnspentOutputs)
		table.AddRow("", "BALANCE", fmt.Sprintf("%f XMR",
			float64(saddr.Balance)/float64(constant.XMR)))

		if saddr.BlocksToUnlock > 0 {
			table.AddRow("", "BLOCKS TO UNLOCK", saddr.BlocksToUnlock)
			table.AddRow("", "TIME TO UNLOCK", saddr.TimeToUnlock)
			table.AddRow("", "UNLOCKED BALANCE", fmt.Sprintf("%f XMR",
				float64(saddr.UnlockedBalance)/float64(constant.XMR)))
		}
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getBalanceCommand{}).Cmd())
}
