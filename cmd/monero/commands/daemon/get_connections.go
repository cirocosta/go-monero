package daemon

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getConnectionsCommand struct {
	JSON bool
}

func (c *getConnectionsCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-connections",
		Short: "information about incoming and outgoing connections.",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getConnectionsCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetConnections(ctx)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	c.pretty(resp)
	return nil
}

func (c *getConnectionsCommand) pretty(v *daemon.GetConnectionsResult) {
	table := display.NewTable()

	table.AddRow("ADDR", "IN", "STATE", "TIME", "RECV (kB)", "SEND (kB)")

	for _, connection := range v.Connections {
		table.AddRow(
			connection.Address,
			connection.Incoming,
			connection.State,
			time.Duration(connection.LiveTime)*time.Second,
			connection.RecvCount/1024,
			connection.SendCount/1024,
		)
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getConnectionsCommand{}).Cmd())
}
