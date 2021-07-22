package daemon

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getTransactionPoolStatsCommand struct {
	JSON bool
}

func (c *getTransactionPoolStatsCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-transaction-pool-stats",
		Short: "statistics about the transaction pool",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(
		&c.JSON,
		"json",
		false,
		"whether or not to output the result as json",
	)

	return cmd
}

func (c *getTransactionPoolStatsCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetTransactionPoolStats(ctx)
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
func (c *getTransactionPoolStatsCommand) pretty(v *daemon.GetTransactionPoolStatsResult) {
	table := display.NewTable()

	table.AddRow("Bytes Max:", v.PoolStats.BytesMax)
	table.AddRow("Bytes Med:", v.PoolStats.BytesMed)
	table.AddRow("Bytes Min:", v.PoolStats.BytesMin)
	table.AddRow("Bytes Total:", v.PoolStats.BytesTotal)
	table.AddRow("Fee Total:", v.PoolStats.FeeTotal)
	table.AddRow("Histogram 98pct:", v.PoolStats.Histo98Pc)
	table.AddRow("Txns in Pool for Longer than 10m:", v.PoolStats.Num10M)
	table.AddRow("Double Spends:", v.PoolStats.NumDoubleSpends)
	table.AddRow("Failing Transactions:", v.PoolStats.NumFailing)
	table.AddRow("Not Relayed:", v.PoolStats.NumNotRelayed)
	table.AddRow("Oldest:", humanize.Time(time.Unix(v.PoolStats.Oldest, 0)))
	table.AddRow("Txns Total:", v.PoolStats.TxsTotal)

	table.AddRow("")
	table.AddRow("BYTES", "TXNS")
	for _, h := range v.PoolStats.Histo {
		if h.Bytes == 0 {
			continue
		}
		table.AddRow(h.Bytes, h.Txs)
	}

	fmt.Println(table)
}

func init() {
	RootCommand.AddCommand((&getTransactionPoolStatsCommand{}).Cmd())
}
