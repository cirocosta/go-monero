package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/display"
	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/constant"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

type getBlockCommand struct {
	Height    uint64
	Hash      string
	BlockJSON bool
	JSON      bool

	client *daemon.Client
}

func (c *getBlockCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-block",
		Short: "full block information by either block height or hash",
		RunE:  c.RunE,
	}

	cmd.Flags().Uint64Var(&c.Height, "height",
		0, "height of the block to retrieve the information of")

	cmd.Flags().StringVar(&c.Hash, "hash",
		"", "block hash to retrieve the information of")

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")
	cmd.Flags().BoolVar(&c.BlockJSON, "block-json",
		false, "display just the block json (from the `json` field)")

	return cmd
}

func (c *getBlockCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	if c.Hash == "" && c.Height == 0 {
		return fmt.Errorf("hash or height must be set")
	}

	resp, err := client.GetBlock(ctx, daemon.GetBlockRequestParameters{
		Hash:   c.Hash,
		Height: c.Height,
	})
	if err != nil {
		return fmt.Errorf("get block: %w", err)
	}

	if c.JSON {
		if !c.BlockJSON {
			return display.JSON(resp)
		}

		inner, err := resp.InnerJSON()
		if err != nil {
			return fmt.Errorf("inner json: %w", err)
		}

		return display.JSON(inner)
	}

	c.client = client
	return c.pretty(ctx, resp)
}

func (c *getBlockCommand) pretty(ctx context.Context, v *daemon.GetBlockResult) error {
	table := display.NewTable()

	blockDetails, err := v.InnerJSON()
	if err != nil {
		return fmt.Errorf("inner json: %w", err)
	}

	table.AddRow("Hash:", v.BlockHeader.Hash)
	table.AddRow("Height:", v.BlockHeader.Height)
	table.AddRow("Age:", humanize.Time(time.Unix(v.BlockHeader.Timestamp, 0)))
	table.AddRow("Timestamp:", time.Unix(v.BlockHeader.Timestamp, 0))
	table.AddRow("Size:", humanize.IBytes(v.BlockHeader.BlockSize))
	table.AddRow("Reward:", fmt.Sprintf("%f XMR", float64(v.BlockHeader.Reward)/constant.XMR))
	table.AddRow("Version:", fmt.Sprintf("%d.%d", blockDetails.MajorVersion, blockDetails.MinorVersion))
	table.AddRow("Previous Block:", blockDetails.PrevID)
	table.AddRow("Nonce:", blockDetails.Nonce)
	table.AddRow("Miner TXN Hash:", v.MinerTxHash)
	fmt.Println(table)
	fmt.Println("")

	txnsResult, err := c.client.GetTransactions(ctx, blockDetails.TxHashes)
	if err != nil {
		return fmt.Errorf("get txns: %w", err)
	}

	table = display.NewTable()
	table.AddRow("HASH", "FEE (µɱ)", "FEE (µɱ per kB)", "IN/OUT", "SIZE")
	for _, txn := range txnsResult.Txs {
		txnDetails := &daemon.TransactionJSON{}
		if err := json.Unmarshal([]byte(txn.AsJSON), txnDetails); err != nil {
			return fmt.Errorf("unsmarshal txjson: %w", err)
		}

		fee := float64(txnDetails.RctSignatures.Txnfee)
		size := len(txn.AsHex) / 2

		table.AddRow(
			txn.TxHash,
			fee/constant.MicroXMR,
			fmt.Sprintf("%6.1f", (fee/constant.MicroXMR)/(float64(size)/1024)),
			fmt.Sprintf("%d/%d", len(txnDetails.Vin), len(txnDetails.Vout)),
			humanize.IBytes(uint64(size)),
		)
	}

	fmt.Println(table)
	return nil
}

func init() {
	RootCommand.AddCommand((&getBlockCommand{}).Cmd())
}
