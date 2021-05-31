# go-monero

A Go library (and CLI) for interacting with Monero daemons via RPC or the P2P
network free of CGO, either on clearnet or not.

- want/need help? reach out to `utxobr` on
  https://matrix.to/#/#monero-community:matrix.org (`#monero-dev`)

## Quick start

### Library

```console
$ go get -u -v github.com/cirocosta/go-monero
```

`go-monero` exposes two high-level packages: `levin` and `daemonrpc`.

The first (`levin`) is used for interacting with the p2p network via plain TCP
(optionally, Tor and I2P can also be used via socks5 proxy - see options). For
instance, to reach out to a node (of a particular address `addr`) and grab its
list of connected peers (information that comes out of the initial handshake):

```golang
import (
        "fmt"
        "context"

        "github.com/cirocosta/go-monero/pkg/levin
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

The second (`daemonrpc`), is used to communicate with `monerod` via its HTTP
endpoints. Note that not all endpoints/fields are exposed on a given port - if
it's being served in a restricted manner, you'll have access to less endpoints
than you see in the documentation
(https://www.getmonero.org/resources/developer-guides/daemon-rpc.html)

For instance:

```go
import (
        "fmt"
        "context"

        "github.com/cirocosta/go-monero/pkg/daemonrpc"
)

func ShowBlockHeight (ctx context.Context, addr string) error {
	client, err := daemonrpc.NewClient(addr)
	if err != nil {
		return fmt.Errorf("new client for '%s': %w", addr, err)
	}

	resp, err := client.GetBlockCount()
	if err != nil {
		return fmt.Errorf("get block count: %w", err)
	}

        fmt.Println(resp.Count)
	return nil
}
```


### CLI

Under `cmd/monero` you'll find a command line interface that exposes most of
the functionality that the library provides.

```console
$ GO111MODULE=on go get github.com/cirocosta/go-monero/cmd/monero

$ monero --help
Usage:
  monero [OPTIONS] <command>

Application Options:
  -v, --verbose  dump http requests and responses to stderr [$MONEROD_VERBOSE]
  -a, --address= RPC server address [$MONEROD_ADDRESS]

Help Options:
  -h, --help     Show this help message

Available commands:
  crawl                       Crawl over the network to find all peers

  p2p-peer-list               Find out the list of local peers known by a node

  get-alternate-chains        Get alternate chains
  get-bans                    Get bans
  get-block                   Get block
  get-block-count             Get the block count
  get-block-template          Get a block template on which mining a new block
  get-coinbase-tx-sum         Get the coinbase amount and the fees amount for n last blocks starting at particular height
  get-connections             Retrieve information about incoming and outgoing connections to your node (restricted)
  get-fee-estimate            Gives an estimation on fees per byte
  get-info                    Retrieve general information about the state of your node and the network. (restricted)
  get-last-block-header       Get the header of the last block
  get-peer-list               Get peer list
  get-transaction-pool        Get all transactions in the pool
  get-transaction-pool-stats  Get the transaction pool statistics
  get-transactions            Retrieve transactions
  hard-fork-info              Get hard fork info
  on-get-block-hash           Look up a block's hash by its height
  sync-info                   Get synchronisation information (restricted)
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
