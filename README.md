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

        STEALTH ADDR                                                            AMOUNT  AMOUNT IDX
0       574e3a3dda7cde249a81c0a0637ce78999e114dd478479e3e04680ffdcd34c97        ?       35307867
1       7d3e93ecfa2bd8c27f7a2bbd57488c1121924e58591b32e99be42da6ae249e9c        ?       35307868


Input Key Image:        2cd588cbc9214e1b683a8c3f0d7c25b8e61779806051356dee1fac6cf42c0c7e

        RING MEMBER                                                             TXID                                                                    BLK     AGE
0       ef729a26d047508fb1b37a553053c8959eac6191bbc5b7c2845aa608d7f48d08        1ed88a49d533e01a222ca12c2fa74fa1797a667509e1e1b6e754cb90944639b0        2342514 2 months ago
1       33351ad5aed5df1c4b38d995c9461112899223bee8a156a77eebad43a2f62c21        e0782b630670dba4062aec0a1f9678181e8321dcaf22638cf7a6033957d5e4f9        2382635 1 month ago
2       df5b27b873dc93e4237031cc598062db86e011e3eec021dd57043076d4e4bfca        a96889571eab740fffd242ac5742b1af83471a5afc0b6ee8e163c4f4fa29574e        2397359 1 week ago
3       f30b7e8a36b38d95dcbbbdf7f291f452d540d02d3d87c791a106682462ef62b7        a3eede407f1069c959c75295a0c0dc15ffae55c5b17ebfffa0ac70f9304fd943        2401095 1 week ago
4       f684f88cd52ef0c39b3466c68eef1c71590ccf42f0936ae2996b1bd0f20dddcc        de3a7e7b8ced8d272cd039c8d12c81f5297766f2cdfd86910b634d47f6faf3c3        2401120 1 week ago
5       96b6addb19dfa8ea78b167dd46acd04d90e4628c17f9cb21ed82f6a0e2ce6c55        13c9063f09e0559eb354a674fcff8dad2288cbf198b64c2d32a5ab8315039dc5        2402375 6 days ago
6       cd573966c65055537747004e702f6cf7d7d32a50a1eed4e8291eca102d0072b4        4639ecf7293c8fd58cf2b34acaf9dd147e53018270b2cd2ddb4e719b28e34e9b        2403122 5 days ago
7       b9c6c95c099426624936d142c2e0a09e641cad0295d468e0bca487502abafe51        04800ff34ee29f63a583d84567f90a3e8ee5f018c17ca4ccbb6afb3305b3bdd9        2404029 3 days ago
8       c4b4d2549c9a73115747ea664d85e4becb502287224eb5bd7b5b41394cd67976        6ebad604e6d81afba32773401f3006f74a51c69e50a98b55e5bd0f0910f08dea        2405504 1 day ago
9       36516e699ceffb665b1cb2bae1e116da03927718a4b391ecaf9f216661494138        e996d5b951ed1484b10d025f96404f54140164d225c044cb2fcb67458b022cc0        2406458 9 hours ago
10      91ab8220d7d87e465fd53521f2714388074b66bc0ae3d2b7ec89f949ca6d6718        9dd40e49611f14380b53202d4d669ce0a141ac308204e9d3e917cb2453f88716        2406541 6 hours ago
```


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
