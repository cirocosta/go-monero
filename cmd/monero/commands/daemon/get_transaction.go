package daemon

import (
	"context"
	"encoding/hex"
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

type getTransactionCommand struct {
	Txn    string
	Unwrap bool
	JSON   bool

	client *daemon.Client
}

func (c *getTransactionCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-transaction",
		Short: "lookup a transaction, in the pool or not",
		RunE:  c.RunE,
	}

	cmd.Flags().StringVar(&c.Txn, "txn",
		"", "hash of a transaction to lookup")
	cmd.MarkFlagRequired("txn")

	cmd.Flags().BoolVar(&c.JSON, "json",
		false, "whether or not to output the result as json")
	cmd.Flags().BoolVar(&c.Unwrap, "unwrap",
		false, "whether or not to unwrap the json representation of the transaction")

	return cmd
}

func (c *getTransactionCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.Client()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.GetTransactions(ctx, []string{c.Txn})
	if err != nil {
		return fmt.Errorf("get transactions: %w", err)
	}

	if c.JSON {
		if !c.Unwrap {
			return display.JSON(resp)
		}

		txns, err := resp.GetTransactions()
		if err != nil {
			return fmt.Errorf("resp get txns: %w", err)
		}

		return display.JSON(txns)
	}

	c.client = client
	return c.pretty(ctx, resp)
}

func (c *getTransactionCommand) pretty(ctx context.Context, v *daemon.GetTransactionsResult) error {
	if len(v.Txs) == 0 {
		return nil
	}

	txn := v.Txs[0]

	txnDetails := &daemon.TransactionJSON{}
	if err := json.Unmarshal([]byte(txn.AsJSON), txnDetails); err != nil {
		return fmt.Errorf("unsmarshal txjson: %w", err)
	}

	if err := c.prettyHeader(ctx, txn, txnDetails); err != nil {
		return err
	}
	if err := c.prettyOutputs(ctx, txn, txnDetails); err != nil {
		return err
	}
	if err := c.prettyInputs(ctx, txn, txnDetails); err != nil {
		return err
	}

	return nil
}

func (c *getTransactionCommand) prettyHeader(
	ctx context.Context,
	txn daemon.GetTransactionsResultTransaction,
	txnDetails *daemon.TransactionJSON,
) error {
	table := display.NewTable()

	confirmations := uint64(0)
	fee := float64(txnDetails.RctSignatures.Txnfee)
	size := len(txn.AsHex) / 2

	table.AddRow("Hash:", txn.TxHash)
	table.AddRow("Fee (µɱ):", fee/constant.MicroXMR)
	table.AddRow("Fee per kB (µɱ):", (fee/constant.MicroXMR)/(float64(size)/1024))
	table.AddRow("In/Out:", fmt.Sprintf("%d/%d", len(txnDetails.Vin), len(txnDetails.Vout)))
	table.AddRow("Size:", humanize.IBytes(uint64(len(txn.AsHex))/2))
	table.AddRow("Public Key:", hex.EncodeToString(txnDetails.Extra[1:33]))

	if txn.InPool == false {
		table.AddRow("Age:", humanize.Time(time.Unix(txn.BlockTimestamp, 0)))
		table.AddRow("Block:", txn.BlockHeight)

		heightResp, err := c.client.GetHeight(ctx)
		if err != nil {
			return fmt.Errorf("get block count: %w", err)
		}

		confirmations = heightResp.Height - txn.BlockHeight
		table.AddRow("Confirmations:", confirmations)
	} else {
		table.AddRow("Confirmations:", 0)
	}

	fmt.Println(table)
	fmt.Println("")

	return nil
}

func (c *getTransactionCommand) prettyOutputs(
	ctx context.Context,
	txn daemon.GetTransactionsResultTransaction,
	txnDetails *daemon.TransactionJSON,
) error {
	table := display.NewTable()
	table.AddRow("OUTPUTS")
	table.AddRow("", "STEALTH ADDR", "AMOUNT", "AMOUNT IDX")
	for idx, vout := range txnDetails.Vout {
		table.AddRow(
			idx,
			vout.Target.Key,
			vout.Amount,
			txn.OutputIndices[idx],
		)
	}

	fmt.Println(table)
	fmt.Println("")

	return nil
}

func decodeOffsets(offsets []uint) []uint {
	accum := uint(0)
	res := make([]uint, len(offsets))

	for idx, offset := range offsets {
		accum += offset
		res[idx] = accum
	}

	return res
}

func (c *getTransactionCommand) prettyInputs(
	ctx context.Context,
	txn daemon.GetTransactionsResultTransaction,
	txnDetails *daemon.TransactionJSON,
) error {
	for _, vin := range txnDetails.Vin {
		outsResp, err := c.client.GetOuts(ctx, decodeOffsets(vin.Key.KeyOffsets), true)
		if err != nil {
			return fmt.Errorf("outs: %w", err)
		}

		table := display.NewTable()
		table.AddRow("Input Key Image:", vin.Key.KImage)
		fmt.Println(table)

		table = display.NewTable()
		table.AddRow("", "RING MEMBER", "TXID", "BLK", "AGE")
		for idx, out := range outsResp.Outs {
			blockHeaderResp, err := c.client.GetBlockHeaderByHeight(ctx, out.Height)
			if err != nil {
				return fmt.Errorf("get block header by height %d: %w", out.Height, err)
			}

			table.AddRow(idx, out.Key, out.Txid, out.Height,
				humanize.Time(time.Unix(blockHeaderResp.BlockHeader.Timestamp, 0)),
			)

		}
		fmt.Println(table)
		fmt.Println()
	}

	return nil
}

func init() {
	RootCommand.AddCommand((&getTransactionCommand{}).Cmd())
}
