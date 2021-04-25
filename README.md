# go-monero

A Go library (and CLI) for interacting with Monero daemons via RPC or the P2P
network free of CGO, either on clearnet or not.

- want/need help? reach out to `utxobr` on
  https://matrix.to/#/#monero-community:matrix.org (`#monero-dev`)

## Quick start

### library

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
	client, err := levin.NewClient(ctx, addr)
	if err != nil {
		return fmt.Errorf("new client '%s': %w", addr, err)
	}

	defer client.Close()

	pl, err := client.Handshake(ctx)
	if err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

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


### cli

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
  crawl                 Crawl over the network to find all peers

  p2p-peer-list         Find out the list of local peers known by a node

  get-block-count       Get the block count
  get-block-template    Get a block template on which mining a new block
  get-coinbase-tx-sum   Get the coinbase amount and the fees amount for n last blocks starting at particular height
  get-connections       Retrieve information about incoming and outgoing connections to your node (restricted)
  get-fee-estimate      Gives an estimation on fees per byte
  get-info              Retrieve general information about the state of your node and the network. (restricted)
  get-transaction-pool  Get all transactions in the pool
  on-get-block-hash     Look up a block's hash by its height
  sync-info             Get synchronisation information (restricted)
```

## TODO

### levin

- [x] header
- [x] payload
  - [x] serialization
  - [x] deserialization
- [ ] properly typed structures

### daemon

#### non-jsonrpc

- [ ] `flush_txpool`
- [ ] `get_alternate_chains`
- [ ] `get_bans`
- [ ] `get_block_header_by_hash`
- [ ] `get_block_header_by_height`
- [ ] `get_block_headers_range`
- [ ] `get_block`
- [ ] `get_last_block_header`
- [ ] `get_output_distribution`
- [ ] `get_output_histogram`
- [ ] `get_txpool_backlog`
- [ ] `get_version`
- [ ] `hard_fork_info`
- [ ] `relay_tx`
- [ ] `set_bans`
- [ ] `submit_block`
- [x] `get_block_count`
- [x] `get_block_template`
- [x] `get_coinbase_tx_sum`
- [x] `get_connections`
- [x] `get_fee_estimate`
- [x] `get_info`
- [x] `on_get_block_hash`
- [x] `sync_info`

#### json rpc

- [ ] `/get_alt_blocks_hashes`
- [ ] `/get_blocks.bin`
- [ ] `/get_blocks_by_height.bin`
- [ ] `/get_hashes.bin`
- [ ] `/get_height`
- [ ] `/get_info (not JSON)`
- [ ] `/get_limit`
- [ ] `/get_o_indexes.bin`
- [ ] `/get_outs.bin`
- [ ] `/get_outs`
- [ ] `/get_peer_list`
- [ ] `/get_transaction_pool_hashes.bin`
- [ ] `/get_transaction_pool_stats`
- [ ] `/get_transactions`
- [ ] `/in_peers`
- [ ] `/is_key_image_spent`
- [ ] `/mining_status`
- [ ] `/out_peers`
- [ ] `/save_bc`
- [ ] `/send_raw_transaction`
- [ ] `/set_limit`
- [ ] `/set_log_categories`
- [ ] `/set_log_hash_rate`
- [ ] `/set_log_level`
- [ ] `/start_mining`
- [ ] `/start_save_graph`
- [ ] `/stop_daemon`
- [ ] `/stop_mining`
- [ ] `/stop_save_graph`
- [ ] `/update`
- [x] `/get_transaction_pool`


## thanks

- `#monero-dev` (https://matrix.to/#/#freenode_#monero-dev:matrix.org)
- https://github.com/cdiv1e12/py-levin
- https://github.com/cryptonotefoundation/cryptonote
- https://github.com/LeTurt/turtlegod
