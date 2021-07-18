package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getLastBlockHeaderCommand struct {
	JSON bool
}

func (c *getLastBlockHeaderCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-last-block-header",
		Short: "header of the last block.",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getLastBlockHeaderCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetLastBlockHeader(ctx)
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
func (c *getLastBlockHeaderCommand) pretty(v *daemon.GetLastBlockHeaderResult) {
	table := display.NewTable()

	prettyBlockHeader(table, v.BlockHeader)

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getLastBlockHeaderCommand{}).Cmd())
}
