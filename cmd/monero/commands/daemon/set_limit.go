package daemon

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type setLimitCommand struct {
	LimitUp   uint64
	LimitDown uint64

	JSON bool
}

func (c *setLimitCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-limit",
		Short: "configure bandwidth throttling",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().Uint64Var(&c.LimitUp, "up",
		0, "max upload bandwidth (in kB/s)")
	cmd.Flags().Uint64Var(&c.LimitDown, "down",
		0, "max download bandwidth (in kB/s)")

	return cmd
}

func (c *setLimitCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	params := daemon.SetLimitRequestParameters{
		LimitUp:   c.LimitUp,
		LimitDown: c.LimitDown,
	}
	resp, err := client.SetLimit(ctx, params)
	if err != nil {
		return fmt.Errorf("set limit: %w", err)
	}
	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *setLimitCommand) pretty(v *daemon.SetLimitResult) {
	table := display.NewTable()
	table.AddRow("Status:", v.Status)
	table.AddRow("Limit Up:", humanize.Bytes(v.LimitUp*1024)+"/s")
	table.AddRow("Limit Down:", humanize.Bytes(v.LimitDown*1024)+"/s")
	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&setLimitCommand{}).Cmd())
}
