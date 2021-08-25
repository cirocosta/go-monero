package monero

import (
	"encoding/binary"
	"hash/crc32"

	"github.com/paxos-bankchain/moneroutil"
)

// KeySize denotes the size of the low-level public and private keys.
//
const KeySize = 32

// SeedOption describes the type of functional options that can be provided to
// the constructor to override default settings.
//
type SeedOption func(s *Seed)

// WithNetwork overrides the default network set for Seed instances.
//
func WithNetwork(n Network) SeedOption {
	return func(s *Seed) {
		s.network = n
	}
}

// Seed encapsulates funcionality that arised from the knowledge of a private
// spend key.
//
// Differently from Bitcoin, Monero users posses two sets of private and public
// keys:
//		  private | public
//                ------- | ------
//	          ks      | Ks		 spend
//	          kv      | Kv		 view
//
// From the private spend key, a private view key is derived. Of each of them,
// a corresponding public key is formed.
//
type Seed struct {
	privateSpendKey []byte
	network         Network

	privateViewKey []byte
	publicSpendKey []byte
	publicViewKey  []byte
}

// NewSeed instantiates a new Seed instance based on a byte array that contains
// the secret number that derives all three other keys.
//
// As mentioned before, although here we store two private keys, ultimately,
// this `privateSpendKey` (an immense number "impossible" to guess that here is
// represented as a 32-byte array) is the only one key that is _really_
// important to safe-guard as the other one is derived out of it.
//
func NewSeed(privateSpendKey []byte, opts ...SeedOption) *Seed {
	s := &Seed{
		privateSpendKey: privateSpendKey,
		network:         NetworkMainnet,

		privateViewKey: make([]byte, KeySize),
		publicSpendKey: make([]byte, KeySize),
		publicViewKey:  make([]byte, KeySize),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.deriveKeys()

	return s
}

// deriveKeys takes the private spend key and, out of it, derives all the other
// three keys:
//
// 	- private view
// 	- public spend
// 	- public view
//
func (s *Seed) deriveKeys() {
	moneroutil.ScReduce32((*moneroutil.Key)(s.privateSpendKey))

	s.privateViewKey = keccak256(s.privateSpendKey)
	s.publicSpendKey = publicKeyFromPrivateKey(s.privateSpendKey)
	s.publicViewKey = publicKeyFromPrivateKey(s.privateViewKey)
}

// Mnemonic converts a private spend key (that gigantic number represented as a
// 32-byte array) to a 25-element list of carefully chosen words in a
// particular language (a 1626 dictionary whose words have some nice
// properties).
//
// The first 24 words correspond to the key, with the 25-th being a CRC32
// checksum used for error-checking.
//
// ps.: in this implementation, only support for the English wordlist is provided.
//
func (s *Seed) Mnemonic() []string {
	mnemonic := make([]string, 25)
	wordlistSize := uint32(len(WordlistEnglish))

	for i := 0; i < KeySize; i += 4 {
		x := binary.LittleEndian.Uint32(s.privateSpendKey[i : i+4])

		w1 := x % wordlistSize
		w2 := (x/wordlistSize + w1) % wordlistSize
		w3 := (x/wordlistSize/wordlistSize + w2) % wordlistSize

		mnemonic[i/4*3] = WordlistEnglish[w1]
		mnemonic[i/4*3+1] = WordlistEnglish[w2]
		mnemonic[i/4*3+2] = WordlistEnglish[w3]
	}

	hash := crc32.NewIEEE()
	for _, word := range mnemonic[:24] {
		rrune := string([]rune(word)[:3])
		hash.Write([]byte(rrune))
	}

	sum := hash.Sum32()
	idx := sum % 24
	mnemonic[24] = mnemonic[idx]

	return mnemonic
}

// PrimaryAddress gives the base58-formatted representation of the primary
// address of this seed.
//
func (s *Seed) PrimaryAddress() string {
	hash := keccak256(
		s.network.PublicAddressBase58Prefix(),
		s.publicSpendKey,
		s.publicViewKey,
	)

	return moneroutil.EncodeMoneroBase58(
		s.network.PublicAddressBase58Prefix(),
		s.publicSpendKey,
		s.publicViewKey,
		hash[:4],
	)
}

func (s *Seed) PrivateSpendKey() []byte {
	return s.privateSpendKey
}

func (s *Seed) PrivateViewKey() []byte {
	return s.privateViewKey
}

func (s *Seed) PublicSpendKey() []byte {
	return s.publicSpendKey
}

func (s *Seed) PublicViewKey() []byte {
	return s.publicViewKey
}
