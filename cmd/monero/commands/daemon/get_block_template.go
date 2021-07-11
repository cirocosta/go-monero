package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
)

type getBlockTemplateCommand struct {
	ExtraNonce    string
	PreviousBlock string
	ReserveSize   uint
	WalletAddress string
}

func (c *getBlockTemplateCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-block-template",
		Short: "generate a block template for mining a new block",
		RunE:  c.RunE,
	}

	cmd.Flags().StringVar(&c.WalletAddress, "wallet-address",
		"", "address of the wallet to receive coinbase transactions if block is successfully mined")
	cmd.MarkFlagRequired("wallet-address")

	cmd.Flags().UintVar(&c.ReserveSize, "reserve-size",
		0, "reserve size")

	cmd.Flags().StringVar(&c.PreviousBlock, "previous-block",
		"", "previous block")

	cmd.Flags().StringVar(&c.ExtraNonce, "extra-nonce",
		"", "extra nonce")

	return cmd
}

func (c *getBlockTemplateCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetBlockTemplate(ctx, c.WalletAddress, c.ReserveSize)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	return display.JSON(resp)
}

func init() {
	RootCommand.AddCommand((&getBlockTemplateCommand{}).Cmd())
}
