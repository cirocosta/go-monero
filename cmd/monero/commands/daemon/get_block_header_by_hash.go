package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getBlockHeaderByHashCommand struct {
	Hashes []string
	Unwrap bool

	JSON bool
}

func (c *getBlockHeaderByHashCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-block-header-by-hash",
		Short: "retrieve block(s) header(s) by hash",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().StringArrayVar(&c.Hashes, "hash",
		[]string{}, "hash of the block to get the header of")
	cmd.MarkFlagRequired("hash")

	return cmd
}

func (c *getBlockHeaderByHashCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetBlockHeaderByHash(ctx, daemon.GetBlockHeaderByHashRequestParameters{
		Hashes: c.Hashes,
	})
	if err != nil {
		return fmt.Errorf("get block: %w", err)
	}

	c.pretty(resp)
	return nil
}

func (c *getBlockHeaderByHashCommand) pretty(v *daemon.GetBlockHeaderByHashResult) {
	table := display.NewTable()

	for _, blockHeader := range v.BlockHeaders {
		prettyBlockHeader(table, blockHeader)
		table.AddRow("")
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getBlockHeaderByHashCommand{}).Cmd())
}
