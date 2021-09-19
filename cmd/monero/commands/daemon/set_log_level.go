package daemon

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type setLogLevelCommand struct {
	Level int8

	JSON bool
}

func (c *setLogLevelCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-log-level",
		Short: "configure the daemon log level",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().Int8Var(&c.Level, "level",
		0, "daemon log  level (0-4, from less to more verbose")
	_ = cmd.MarkFlagRequired("level")

	return cmd
}

func (c *setLogLevelCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	params := daemon.SetLogLevelRequestParameters{
		Level: c.Level,
	}
	resp, err := client.SetLogLevel(ctx, params)
	if err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *setLogLevelCommand) pretty(v *daemon.SetLogLevelResult) {
	table := display.NewTable()
	table.AddRow("Status:", v.Status)
	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&setLogLevelCommand{}).Cmd())
}
