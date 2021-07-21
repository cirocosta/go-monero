package wallet

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cirocosta/go-monero/cmd/monero/options"
	"github.com/cirocosta/go-monero/pkg/rpc/wallet"
)

type autoRefreshCommand struct {
	Enable bool
	Period time.Duration
}

func (c *autoRefreshCommand) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auto-refresh",
		Short: "configure wallet-rpc to automatically refresh or not",
		RunE:  c.RunE,
	}

	cmd.Flags().BoolVar(&c.Enable, "enable",
		true, "enable the ability to auto-refresh")

	cmd.Flags().DurationVar(&c.Period, "period",
		90*time.Second, "how often to automatically refresh")

	return cmd
}

func (c *autoRefreshCommand) RunE(_ *cobra.Command, _ []string) error {
	ctx, cancel := options.RootOptions.Context()
	defer cancel()

	client, err := options.RootOptions.WalletClient()
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	resp, err := client.AutoRefresh(ctx, c.Enable, int64(c.Period.Seconds()))
	if err != nil {
		return fmt.Errorf("auto refresh: %w", err)
	}

	c.pretty(resp)
	return nil
}

// nolint:forbidigo
func (c *autoRefreshCommand) pretty(v *wallet.AutoRefreshResult) {
	fmt.Println("OK")
}

func init() {
	RootCommand.AddCommand((&autoRefreshCommand{}).Cmd())
}
