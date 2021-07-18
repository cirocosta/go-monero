package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getBlockHeaderCommand struct {
	Hashes []string
	Height uint64
	Unwrap bool

	JSON bool
}

func (c *getBlockHeaderCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-block-header",
		Short: "retrieve block(s) header(s) by hash",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().Uint64Var(&c.Height, "height",
		0, "height of a block to fetch")
	cmd.Flags().StringArrayVar(&c.Hashes, "hash",
		[]string{}, "hash of the block to get the header of")

	return cmd
}

func (c *getBlockHeaderCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	if len(c.Hashes) > 0 {
		resp, err := client.GetBlockHeaderByHash(ctx, c.Hashes)
		if err != nil {
			return fmt.Errorf("get block header by hash: %w", err)
		}

		c.pretty(resp.BlockHeaders)

		return nil
	}

	resp, err := client.GetBlockHeaderByHeight(ctx, c.Height)
	if err != nil {
		return fmt.Errorf("get block header by height: %w", err)
	}

	c.pretty([]daemon.BlockHeader{resp.BlockHeader})
	return nil
}

// nolint:forbidigo
func (c *getBlockHeaderCommand) pretty(blockHeaders []daemon.BlockHeader) {
	table := display.NewTable()

	for _, blockHeader := range blockHeaders {
		prettyBlockHeader(table, blockHeader)
		table.AddRow("")
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getBlockHeaderCommand{}).Cmd())
}
