package monero

import (
	"github.com/paxos-bankchain/moneroutil"
	"golang.org/x/crypto/sha3"
)

func keccak256(data ...[]byte) []byte {
	hash := sha3.NewLegacyKeccak256()

	for _, v := range data {
		hash.Write(v)
	}

	return hash.Sum(nil)
}

func publicKeyFromPrivateKey(private []byte) []byte {
	public := make([]byte, KeySize)

	p := new(moneroutil.ExtendedGroupElement)
	moneroutil.GeScalarMultBase(p, (*moneroutil.Key)(private))
	p.ToBytes((*moneroutil.Key)(public))

	return public
}
