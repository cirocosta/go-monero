package daemon

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/constant"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getTransactionPoolCommand struct {
	JSON bool
}

func (c *getTransactionPoolCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "get-transaction-pool",
		Short: "information about valid transactions seen by the " +
			"node but not yet mined into a block, including " +
			"spent key image info for the txpool",
		RunE: c.RunE,
	}

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")

	return cmd
}

func (c *getTransactionPoolCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetTransactionPool(ctx)
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

	if c.JSON {
		return display.JSON(resp)
	}

	return c.pretty(resp)
}

// nolint:forbidigo
func (c *getTransactionPoolCommand) pretty(v *daemon.GetTransactionPoolResult) error {
	table := display.NewTable()
	table.AddRow("Spent Key Images:", len(v.SpentKeyImages))
	fmt.Println(table)
	fmt.Println()

	table = display.NewTable()
	table.AddRow("AGE", "HASH", "FEE (µɱ)", "FEE (µɱ per kB)", "IN/OUT", "SIZE")

	sort.Slice(v.Transactions, func(i, j int) bool {
		return v.Transactions[i].ReceiveTime < v.Transactions[j].ReceiveTime
	})
	for _, txn := range v.Transactions {
		txnDetails := &daemon.TransactionJSON{}
		if err := json.Unmarshal([]byte(txn.TxJSON), txnDetails); err != nil {
			return fmt.Errorf("unsmarshal txjson: %w", err)
		}

		table.AddRow(
			humanize.Time(time.Unix(txn.ReceiveTime, 0)),
			txn.IDHash,
			txn.Fee/constant.MicroXMR,
			fmt.Sprintf("%6.1f", (float64(txn.Fee)/constant.MicroXMR)/(float64(txn.BlobSize)/1024)),
			fmt.Sprintf("%d/%d", len(txnDetails.Vin), len(txnDetails.Vout)),
			humanize.IBytes(txn.BlobSize),
		)
	}

	fmt.Println(table)
	return nil
}

func init() {
	RootCommand.AddCommand((&getTransactionPoolCommand{}).Cmd())
}
