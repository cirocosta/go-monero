package daemon

import (
	"context"
	"fmt"

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
	Last      int64
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

	cmd.Flags().Int64Var(&c.Last, "last",
		-1, "get the last Nth block")

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
	ctx, cancel := options.RootOpts.Context()
	defer cancel()

	client, err := options.RootOpts.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	if c.Last >= 0 {
		lastBlockHeaderResp, err := client.GetLastBlockHeader(ctx)
		if err != nil {
			return fmt.Errorf("get last block header: %w", err)
		}

		c.Height = lastBlockHeaderResp.BlockHeader.Height - uint64(c.Last)
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

// nolint:forbidigo
func (c *getBlockCommand) pretty(ctx context.Context, v *daemon.GetBlockResult) error {
	table := display.NewTable()

	blockDetails, err := v.InnerJSON()
	if err != nil {
		return fmt.Errorf("inner json: %w", err)
	}

	txnsResult, err := c.client.GetTransactions(ctx, blockDetails.TxHashes)
	if err != nil {
		return fmt.Errorf("get txns: %w", err)
	}

	details, err := txnsResult.GetTransactions()
	if err != nil {
		return fmt.Errorf("get transactions: %w", err)
	}

	prettyBlockHeader(table, v.BlockHeader)

	fees := uint64(0)
	for _, d := range details {
		fees += d.RctSignatures.Txnfee
	}

	table.AddRow("Fees:", display.PreciseXMR(fees))
	fmt.Println(table)
	fmt.Println("")

	table = display.NewTable()
	table.AddRow("HASH", "FEE (µɱ)", "FEE (µɱ per kB)", "IN/OUT", "SIZE")
	for idx, txn := range txnsResult.Txs {
		txnDetails := details[idx]

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
