package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/options"
)

type onGetBlockHashCommand struct {
	Height uint64
}

func (c *onGetBlockHashCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "on-get-block-hash",
		Short: "find out block's hash by height",
		RunE:  c.RunE,
	}

	cmd.Flags().Uint64Var(&c.Height, "height",
		0, "block height to find the hash for")
	cmd.MarkFlagRequired("height")

	return cmd
}

func (c *onGetBlockHashCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.OnGetBlockHash(ctx, c.Height)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	fmt.Println(resp)
	return nil
}

func init() {
	RootCommand.AddCommand((&onGetBlockHashCommand{}).Cmd())
}
