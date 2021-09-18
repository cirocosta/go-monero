package daemon

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getLimitCommand struct {
	JSON bool
}

func (c *getLimitCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-limit",
		Short: "retrieve bandwidth throttling information",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getLimitCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetLimit(ctx)
	if err != nil {
		return fmt.Errorf("get limit: %w", err)
	}
	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *getLimitCommand) pretty(v *daemon.GetLimitResult) {
	table := display.NewTable()
	table.AddRow("Status:", v.Status)
	table.AddRow("Limit Up:", humanize.Bytes(v.LimitUp*1024)+"/s")
	table.AddRow("Limit Down:", humanize.Bytes(v.LimitDown*1024)+"/s")
	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getLimitCommand{}).Cmd())
}
