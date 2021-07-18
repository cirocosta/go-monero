package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getAlternateChainsCommand struct {
	JSON bool
}

func (c *getAlternateChainsCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-alternate-chains",
		Short: "display alternative chains as seen by the node",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getAlternateChainsCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetAlternateChains(ctx)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *getAlternateChainsCommand) pretty(v *daemon.GetAlternateChainsResult) {
	table := display.NewTable()

	for _, chain := range v.Chains {
		table.AddRow("Block Hash:", chain.BlockHash)
		table.AddRow("Height:", chain.Height)
		table.AddRow("Main Chain Parent Block:", chain.MainChainParentBlock)
		table.AddRow("Length:", chain.Length)
		table.AddRow("Difficulty:", chain.Difficulty)
		table.AddRow("")
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getAlternateChainsCommand{}).Cmd())
}
