package monero

import "fmt"

// Network denotes a type of Monero network to gather information about.
//
type Network string

const (
	NetworkMainnet   Network = "mainnet"
	NetworkTestnet   Network = "testnet"
	NetworkStagenet  Network = "stagenet"
	NetworkFakechain Network = "fakechain"
)

func (n Network) PublicAddressBase58Prefix() []byte {
	switch n {
	case NetworkMainnet:
		return []byte{18}
	case NetworkTestnet:
		return []byte{53}
	case NetworkStagenet:
		return []byte{24}
	case NetworkFakechain:
		return NetworkMainnet.PublicAddressBase58Prefix()
	}

	panic(fmt.Errorf("'%s' is not a valid netowrk", n))
}
