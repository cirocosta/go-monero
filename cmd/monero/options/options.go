package options

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

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
	address string
	mhttp.ClientConfig
}

// Context generates a new `context.Context` already honouring the deadline
// specified in the options.
//
func (o *options) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), o.RequestTimeout)
}

// initializeFromEnv ensures that any variables not supplied via flags have
// been captures from the set of environment variables.
//
func (o *options) initializeFromEnv() {
	if address := os.Getenv("MONERO_ADDRESS"); address != "" {
		o.address = address
	}
}

// Client instantiates a new daemon RPC client based on the options filled.
//
func (o *options) Client() (*daemon.Client, error) {
	o.initializeFromEnv()

	httpClient, err := mhttp.NewClient(o.ClientConfig)
	if err != nil {
		return nil, fmt.Errorf("new httpclient: %w", err)
	}

	client, err := rpc.NewClient(o.address, rpc.WithHTTPClient(httpClient))
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
	o.initializeFromEnv()

	httpClient, err := mhttp.NewClient(o.ClientConfig)
	if err != nil {
		return nil, fmt.Errorf("new httpclient: %w", err)
	}

	client, err := rpc.NewClient(o.address, rpc.WithHTTPClient(httpClient))
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
	cmd.PersistentFlags().BoolVarP(&RootOptions.Verbose,
		"verbose", "v",
		false,
		"dump http requests and responses to stderr")

	cmd.PersistentFlags().StringVarP(&RootOptions.address,
		"address", "a",
		"http://localhost:18081",
		"full address of the monero node to reach out to [MONERO_ADDRESS]")

	cmd.PersistentFlags().StringVarP(&RootOptions.Username,
		"username", "u",
		"",
		"name of the user to use during rpc auth")

	cmd.PersistentFlags().StringVarP(&RootOptions.Password,
		"password", "p",
		"",
		"password to supply for rpc auth")

	cmd.PersistentFlags().BoolVarP(&RootOptions.TLSSkipVerify,
		"tls-skip-verify", "k",
		false,
		"skip verification of certificate chain and host name")

	cmd.PersistentFlags().StringVar(&RootOptions.TLSClientCert,
		"tls-client-cert",
		"",
		"tls client certificate to use when connecting")

	cmd.PersistentFlags().StringVar(&RootOptions.TLSClientKey,
		"tls-client-key",
		"",
		"tls client key to use when connecting")

	cmd.PersistentFlags().StringVar(&RootOptions.TLSCACert,
		"tls-ca-cert",
		"",
		"certificate authority to load")

	cmd.PersistentFlags().DurationVar(&RootOptions.RequestTimeout,
		"request-timeout",
		1*time.Minute,
		"max wait time until considering the request a failure")
}
