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

### Example

```console
$ monero daemon get-transaction --txn 53c1ef0cff73d12294e3055400826905efc397814aacd7208921a9abfd1f6328
Hash:                   53c1ef0cff73d12294e3055400826905efc397814aacd7208921a9abfd1f6328
Fee (µɱ):               9.13
Fee per kB (µɱ):        6.416691832532602
In/Out:                 1/2
Size:                   1.4 KiB
Public Key:             fbc0ac2c62514f68134f543ce5c4efe51ec55a0952eb188dadf22dccc3c5ffdf
Age:                    20 minutes ago
Block:                  2406735
Confirmations:          4

Outputs

        STEALTH ADDR                                                        AMOUNT IDX
0       574e3a3dda7cde249a81c0a0637ce78999e114dd478479e3e04680ffdcd34c97    35307867
1       7d3e93ecfa2bd8c27f7a2bbd57488c1121924e58591b32e99be42da6ae249e9c    35307868


Input Key Image:        2cd588cbc9214e1b683a8c3f0d7c25b8e61779806051356dee1fac6cf42c0c7e

        RING MEMBER              TXID              BLK        AGE
0       ef729a26d047508fb1       1ed88a49d533e     2342514    2 months ago
1       33351ad5aed5df1c4b       e0782b630670d     2382635    1 month ago
2       df5b27b873dc93e423       a96889571eab7     2397359    1 week ago
3       f30b7e8a36b38d95dc       a3eede407f106     2401095    1 week ago
4       f684f88cd52ef0c39b       de3a7e7b8ced8     2401120    1 week ago
5       96b6addb19dfa8ea78       13c9063f09e05     2402375    6 days ago
6       cd573966c650555377       4639ecf7293c8     2403122    5 days ago
7       b9c6c95c0994266249       04800ff34ee29     2404029    3 days ago
8       c4b4d2549c9a731157       6ebad604e6d81     2405504    1 day ago
9       36516e699ceffb665b       e996d5b951ed1     2406458    9 hours ago
10      91ab8220d7d87e465f       9dd40e49611f1     2406541    6 hours ago
```

_(^ ring member and txid shortened just in this README for brevity sake)_


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
  -p, --password string            password to supply for rpc auth
      --request-timeout duration   max wait time until considering the request a failure (default 1m0s)
      --tls-ca-cert string         certificate authority to load
      --tls-client-cert string     tls client certificate to use when connecting
      --tls-client-key string      tls client key to use when connecting
  -k, --tls-skip-verify            skip verification of certificate chain and host name
  -u, --username string            name of the user to use during rpc auth
  -v, --verbose                    dump http requests and responses to stderr

Use "monero daemon [command] --help" for more information about a command.
```

### Tor support

Nodes reachable only through the Tor network (hidden services) _are_ supported
despite the lack of a specific flag for specifying the proxy address. 

For instance:

```console
$ export HTTP_PROXY="socks5://127.0.0.1:9050" 
$ export MONERO_ADDR=http://rbpgdckle3h3vi4wwwrh75usqtoc5r3alohy7yyx57isynvay63nacyd.onion:18089

$ monero daemon --verbose -a $MONERO_ADDR  get-version --verbose
GET /json_rpc HTTP/1.1
Host: rbpgdckle3h3vi4wwwrh75usqtoc5r3alohy7yyx57isynvay63nacyd.onion:18089
User-Agent: Go-http-client/1.1
Content-Length: 49
Content-Type: application/json
Accept-Encoding: gzip

{"id":"0","jsonrpc":"2.0","method":"get_version"}
200
HTTP/1.1 200 Ok
Content-Length: 150
Accept-Ranges: bytes
Content-Type: application/json
Last-Modified: Sun, 18 Jul 2021 21:10:57 GMT
Server: Epee-based

{
  "id": "0",
  "jsonrpc": "2.0",
  "result": {
    "release": true,
    "status": "OK",
    "untrusted": false,
    "version": 196613
  }
}
Release:        true
Major:          3
Minor:          5
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
