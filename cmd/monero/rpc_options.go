package main

import (
	"context"
	"fmt"
	"time"

	mhttp "github.com/cirocosta/go-monero/pkg/http"
	"github.com/cirocosta/go-monero/pkg/rpc"
)

// global RPC configuration available for all commands.
//
type RPCOptions struct {
	Verbose        bool          `short:"v" env:"MONEROD_VERBOSE" long:"verbose" description:"dump http requests and responses to stderr"`
	RequestTimeout time.Duration `short:"t" env:"MONEROD_TIMEOUT" long:"timeout" description:"request timeout" default:"10s"`
	Address        string        `short:"a" env:"MONEROD_ADDRESS" long:"address" description:"address of the node to target" default:"http://localhost:18081"`
}

// Context generates a new `context.Context` already honouring the deadline
// specified in the options.
//
func (o RPCOptions) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), o.RequestTimeout)
}

// Client instantiates a new RPC client based on the options filled.
//
func (o RPCOptions) Client() (*rpc.Client, error) {
	client, err := rpc.NewClient(o.Address, rpc.WithHTTPClient(mhttp.NewHTTPClient(o.Verbose)))
	if err != nil {
		return nil, fmt.Errorf("new client for '%s': %w", o.Address, err)
	}

	return client, nil
}
