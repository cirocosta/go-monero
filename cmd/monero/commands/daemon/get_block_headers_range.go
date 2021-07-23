package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getBlockHeadersRangeCommand struct {
	Start uint64
	End   uint64
	Last  uint64

	JSON bool
}

func (c *getBlockHeadersRangeCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-block-headers-range",
		Short: "retrieve a range of block headers",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().Uint64Var(&c.Start, "start",
		0, "height of the first block in the range")
	cmd.Flags().Uint64Var(&c.End, "end",
		0, "height the last block in the range")
	cmd.Flags().Uint64Var(&c.Last, "last",
		0, "get the last `n` block headers")

	return cmd
}

func (c *getBlockHeadersRangeCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	start, end := c.Start, c.End
	if c.Last != 0 {
		heightResp, err := client.GetHeight(ctx)
		if err != nil {
			return fmt.Errorf("get height: %w", err)
		}

		end = heightResp.Height - 1
		start = end - c.Last
	}

	resp, err := client.GetBlockHeadersRange(ctx, start, end)
	if err != nil {
		return fmt.Errorf("get block header by height: %w", err)
	}

	c.pretty(resp.Headers)
	return nil
}

// nolint:forbidigo
func (c *getBlockHeadersRangeCommand) pretty(blockHeaders []daemon.BlockHeader) {
	table := display.NewTable()

	for _, blockHeader := range blockHeaders {
		prettyBlockHeader(table, blockHeader)
		table.AddRow("")
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getBlockHeadersRangeCommand{}).Cmd())
}
