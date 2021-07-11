package options

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	mhttp "github.com/cirocosta/go-monero/pkg/http"
	"github.com/cirocosta/go-monero/pkg/rpc"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
	"github.com/cirocosta/go-monero/pkg/rpc/wallet"
)

// RootOptions are global options available to all commands under this package.
//
var RootOptions = &options{}

// options is a set of flags that are shared between all commands in this
// package.
//
type options struct {
	verbose        bool
	requestTimeout time.Duration
	address        string
}

// Context generates a new `context.Context` already honouring the deadline
// specified in the options.
//
func (o *options) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), o.requestTimeout)
}

// Client instantiates a new daemon RPC client based on the options filled.
//
func (o *options) Client() (*daemon.Client, error) {
	client, err := rpc.NewClient(o.address,
		rpc.WithHTTPClient(mhttp.NewHTTPClient(o.verbose)),
	)
	if err != nil {
		return nil, fmt.Errorf("new daemon client for '%s': %w",
			o.address, err,
		)
	}

	return daemon.NewClient(client), nil
}

// WalletClient instantiates a new wallet RPC client based on the options
// filled.
//
func (o *options) WalletClient() (*wallet.Client, error) {
	client, err := rpc.NewClient(o.address,
		rpc.WithHTTPClient(mhttp.NewHTTPClient(o.verbose)),
	)
	if err != nil {
		return nil, fmt.Errorf("new daemon client for '%s': %w",
			o.address, err,
		)
	}

	return wallet.NewClient(client), nil
}

// Bind binds the flags defined by `options` to a `cobra` command so that they
// can be filled either via comand arguments or environment variables.
//
func Bind(cmd *cobra.Command) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MONERO_")

	cmd.PersistentFlags().BoolVarP(
		&RootOptions.verbose,
		"verbose",
		"v",
		false,
		"dump http requests and responses to stderr",
	)
	viper.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))

	cmd.PersistentFlags().StringVarP(
		&RootOptions.address,
		"address",
		"a",
		"http://localhost:18081",
		"full address of the monero node to reach out to",
	)
	viper.BindPFlag("address", cmd.PersistentFlags().Lookup("address"))

	cmd.PersistentFlags().DurationVar(
		&RootOptions.requestTimeout,
		"request-timeout",
		1*time.Minute,
		"how long to wait until considering the request a failure",
	)
	viper.BindPFlag("request-timeout", cmd.PersistentFlags().Lookup("request-timeout"))
}
