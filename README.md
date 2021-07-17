# go-monero

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/cirocosta/go-monero)


A Go library (and CLI) for interacting with Monero's daemon via RPC or the P2P
network, free of CGO, either on clearnet or not.

Support for `monero-wallet-rpc` coming soon.


## Quick start

### Command Line Interface

Under `cmd/monero` you'll find a command line interface that exposes most of
the functionality that the library provides.

You can either install it by using Go building from scratch

```console
$ GO111MODULE=on go get github.com/cirocosta/go-monero/cmd/monero
```

or fetching the binary for your distribution - check out the [releases page](https://github.com/cirocosta/go-monero/releases).

```console
$ monero --help
Daemon, Wallet, and p2p command line monero CLI

Usage:
  monero [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  daemon      execute remote procedure calls against a monero node
  help        Help about any command
  p2p         execute p2p commands against a monero node
  wallet      execute remote procedure calls against a monero wallet rpc server

Flags:
  -h, --help   help for monero

Use "monero [command] --help" for more information about a command.
```


```console
$ monero daemon --help
execute remote procedure calls against a monero node

Usage:
  monero daemon [command]

Available Commands:
  generate-blocks            generate blocks when in regtest mode
  get-alternate-chains       display alternative chains as seen by the node
  get-bans                   all the nodes that have been banned by our node
  get-block                  full block information by either block height or hash
  get-block-count            look up how many blocks are in the longest chain known to the node
  get-block-header           retrieve block(s) header(s) by hash
  get-block-template         generate a block template for mining a new block
  get-coinbase-tx-sum        compute the coinbase amount and the fees amount for n last blocks starting at particular height
  get-connections            information about incoming and outgoing connections.
  get-fee-estimate           estimate fees in atomic units per kB
  get-height                 node's current chain height
  get-info                   general information about the node and the network
  get-last-block-header      header of the last block.
  get-net-stats              networking statistics.
  get-outs                   output details
  get-peer-list              peers lists (white and gray)
  get-public-nodes           all known peers advertising as public nodes
  get-transaction            lookup a transaction, in the pool or not
  get-transaction-pool       information about valid transactions seen by the node but not yet mined into a block, including spent key image info for the txpool
  get-transaction-pool-stats statistics about the transaction pool
  get-version                version of the monero daemon
  hardfork-info              information regarding hard fork voting and readiness.
  on-get-block-hash          find out block's hash by height
  relay-tx                   relay a list of transaction ids
  rpc-access-tracking        statistics about rpc access
  set-bans                   ban another nodes
  sync-info                  daemon's chain synchronization info

Flags:
  -a, --address string             full address of the monero node to reach out to (default "http://localhost:18081")
  -h, --help                       help for daemon
      --request-timeout duration   how long to wait until considering the request a failure (default 1m0s)
  -v, --verbose                    dump http requests and responses to stderr

Use "monero daemon [command] --help" for more information about a command.
```


### Library

To consume `go-monero` as a library for your Go project:

```console
$ go get -u -v github.com/cirocosta/go-monero
```

`go-monero` exposes two high-level packages: `levin` and `rpc`.

The first (`levin`) is used for interacting with the p2p network via plain TCP
(optionally, Tor and I2P can also be used via socks5 proxy - see options). 

For instance, to reach out to a node (of a particular address `addr`) and grab
its list of connected peers (information that comes out of the initial
handshake):

```golang
import (
	"context"
	"fmt"

	"github.com/cirocosta/go-monero/pkg/levin"
)

func ListNodePeers(ctx context.Context, addr string) error {
	// start a client - this will actually establish a TCP `connect()`ion
	// with the other node.
	//
	client, err := levin.NewClient(ctx, addr)
	if err != nil {
		return fmt.Errorf("new client '%s': %w", addr, err)
	}

	// close the connection when done
	//
	defer client.Close()

	// perform the handshake
	//
	pl, err := client.Handshake(ctx)
	if err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	// list the peers reported back (250 max per monero's implementation)
	//
	for addr := range pl.Peers {
		fmt.Println(addr)
	}

	return nil
}
```

The second (`rpc`), is used to communicate with `monerod` via its HTTP
endpoints. Note that not all endpoints/fields are exposed on a given port - if
it's being served in a restricted manner, you'll have access to less endpoints
than you see in the documentation
(https://www.getmonero.org/resources/developer-guides/daemon-rpc.html)

`rpc` itself is subdivided in two other packages: `wallet` and `daemon`, exposing `monero-wallet-rpc` and `monerod` RPCs accordingly.

For instance, to get the the height of the main chain:

```go
package daemon_test

import (
	"context"
	"fmt"

	"github.com/cirocosta/go-monero/pkg/rpc"
	"github.com/cirocosta/go-monero/pkg/rpc/daemon"
)

func ExampleGetHeight() {
	ctx := context.Background()
	addr := "http://localhost:18081"

	// instantiate a generic RPC client
	//
	client, err := rpc.NewClient(addr)
	if err != nil {
		panic(fmt.Errorf("new client for '%s': %w", addr, err))
	}

	// instantiate a daemon-specific client and call the `get_height`
	// remote procedure.
	//
	height, err := daemon.NewClient(client).GetHeight(ctx)
	if err != nil {
		panic(fmt.Errorf("get height: %w", err))
	}

	fmt.Printf("height=%d hash=%s\n", height.Height, height.Hash)
}
```

## License

See [LICENSE](./LICENSE).


## Thanks

Big thanks to the Monero community and other projects around cryptonote:

- `#monero-dev` (https://matrix.to/#/#freenode_#monero-dev:matrix.org)
- https://reddit.com/r/Monero
- https://github.com/cdiv1e12/py-levin
- https://github.com/cryptonotefoundation/cryptonote
- https://github.com/LeTurt/turtlegod


## Donate

![xmr address](./assets/donate.png)

891B5keCnwXN14hA9FoAzGFtaWmcuLjTDT5aRTp65juBLkbNpEhLNfgcBn6aWdGuBqBnSThqMPsGRjWVQadCrhoAT6CnSL3
