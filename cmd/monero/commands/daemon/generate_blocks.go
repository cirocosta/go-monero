package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type generateBlocksCommand struct {
	amountOfBlocks uint64
	walletAddress  string

	JSON bool
}

func (c *generateBlocksCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-blocks",
		Short: "generate blocks when in regtest mode",
		RunE:  c.RunE,
	}

	cmd.Flags().Uint64Var(&c.amountOfBlocks, "amount-of-blocks",
		1, "number of blocks to generate")
	cmd.Flags().StringVar(&c.walletAddress, "wallet-address",
		"", "address submit the block rewards to")
	_ = cmd.MarkFlagRequired("wallet-address")

	return cmd
}

func (c *generateBlocksCommand) RunE(cmd *cobra.Command, args []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GenerateBlocks(ctx, daemon.GenerateBlocksRequestParameters{
		WalletAddress:  c.walletAddress,
		AmountOfBlocks: c.amountOfBlocks,
	})
	if err != nil {
		return fmt.Errorf("generate blocks: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)

	return nil
}

// nolint:forbidigo
func (c *generateBlocksCommand) pretty(v *daemon.GenerateBlocksResult) {
	table := display.NewTable()

	table.AddRow("Final Height:", v.Height)
	for _, block := range v.Blocks {
		table.AddRow("Block:", block)
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&generateBlocksCommand{}).Cmd())
}
