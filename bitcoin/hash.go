package bitcoin

import (
	"crypto/sha256"

	"github.com/alivanz/go-crypto"
)

var BitcoinHasher crypto.Hasher = btchash{}

type btchash struct{}

func (btchash) Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash[:])
	return hash2[:]
}
