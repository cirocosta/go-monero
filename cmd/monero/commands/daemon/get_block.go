package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getBlockCommand struct {
	Height    uint64
	Hash      string
	BlockJSON bool
}

func (c *getBlockCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-block",
		Short: "full block information by either block height or hash",
		RunE:  c.RunE,
	}

	cmd.Flags().Uint64Var(&c.Height, "height",
		0, "height of the block to retrieve the information of")

	cmd.Flags().StringVar(&c.Hash, "hash",
		"", "block hash to retrieve the information of")

	cmd.Flags().BoolVar(&c.BlockJSON, "block-json",
		false, "display just the block json (from the `json` field)")

	return cmd
}

func (c *getBlockCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetBlock(ctx, daemon.GetBlockRequestParameters{
		Hash:   c.Hash,
		Height: c.Height,
	})
	if err != nil {
		return fmt.Errorf("get block: %w", err)
	}

	if !c.BlockJSON {
		return display.JSON(resp)
	}

	inner, err := resp.InnerJSON()
	if err != nil {
		return fmt.Errorf("inner json: %w", err)
	}

	return display.JSON(inner)
}

func init() {
	RootCommand.AddCommand((&getBlockCommand{}).Cmd())
}
