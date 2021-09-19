package daemon

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type setLogCategoriesCommand struct {
	Categories string

	JSON bool
}

func (c *setLogCategoriesCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-log-categories",
		Short: "set the categories for which logs should be emitted",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	cmd.Flags().StringVar(&c.Categories, "categories",
		"", "comma-separated list of <category>:<level> pairs")
	_ = cmd.MarkFlagRequired("categories")

	return cmd
}

func (c *setLogCategoriesCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	params := daemon.SetLogCategoriesRequestParameters{
		Categories: c.Categories,
	}
	resp, err := client.SetLogCategories(ctx, params)
	if err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
	if c.JSON {
		return display.JSON(resp)
	}

	return c.pretty(resp)
}

// nolint:forbidigo
func (c *setLogCategoriesCommand) pretty(v *daemon.SetLogCategoriesResult) error {
	table := display.NewTable()
	table.AddRow("Status:", v.Status)
	table.AddRow("")
	fmt.Println(table)

	table = display.NewTable()
	table.AddRow("CATEGORY", "LEVEL")
	for _, pair := range strings.Split(v.Categories, ",") {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("expected 2 parts, got %d", len(parts))
		}

		table.AddRow(parts[0], parts[1])
	}

	fmt.Println(table)

	return nil
}

func init() {
	RootCommand.AddCommand((&setLogCategoriesCommand{}).Cmd())
}
