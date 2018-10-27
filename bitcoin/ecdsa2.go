package bitcoin

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"

	gcrypto "github.com/alivanz/go-crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

type signer2 struct {
	*ecdsa.PrivateKey
}

func NewWallet2(privkey []byte) (gcrypto.Wallet, error) {
	c := secp256k1.S256()
	D := big.NewInt(0)
	D.SetBytes(privkey)

	pk := new(ecdsa.PrivateKey)
	pk.PublicKey.Curve = c
	pk.D = D
	pk.PublicKey.X, pk.PublicKey.Y = c.ScalarBaseMult(D.Bytes())

	return &signer2{pk}, nil
}

func (x *signer2) PubKey() (*ecdsa.PublicKey, error) {
	return &x.PublicKey, nil
}

func (x *signer2) Sign(msghash []byte) (*big.Int, *big.Int, error) {
	return ecdsa.Sign(rand.Reader, x.PrivateKey, msghash)
}
func (x *signer2) Verify(msghash []byte, r, s *big.Int) bool {
	return ecdsa.Verify(&x.PublicKey, msghash, r, s)
}
