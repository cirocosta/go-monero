package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getVersionCommand struct {
	JSON bool
}

func (c *getVersionCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-version",
		Short: "version of the monero daemon",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getVersionCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetVersion(ctx)
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
func (c *getVersionCommand) pretty(v *daemon.GetVersionResult) {
	table := display.NewTable()

	table.AddRow("Release:", v.Release)
	table.AddRow("Major:", v.Version>>16)
	table.AddRow("Minor:", v.Version&((1<<16)-1))

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getVersionCommand{}).Cmd())
}
