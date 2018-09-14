package bitcoin

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/alivanz/go-utils"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	".."
)

type signer struct {
	*ecdsa.PrivateKey
}

func NewWallet(privkey []byte) (crypto.Wallet, error) {
	c := secp256k1.S256()
	D := big.NewInt(0)
	D.SetBytes(privkey)

	pk := new(ecdsa.PrivateKey)
	pk.PublicKey.Curve = secp256k1.S256()
	pk.D = D
	pk.PublicKey.X, pk.PublicKey.Y = c.ScalarBaseMult(D.Bytes())

	return &signer{pk}, nil
}

func (x *signer) Sign(hash []byte) ([]byte, error) {
	// sign
	r, s, err := ecdsa.Sign(rand.New(rand.NewSource(time.Now().UnixNano())), x.PrivateKey, hash)
	if err != nil {
		return nil, err
	}
	return DERSignature(r, s), nil
}

func (x *signer) PubKey() []byte {
	return append(x.PublicKey.X.Bytes(), x.PublicKey.Y.Bytes()...)
}

func (x *signer) Verify(hash []byte, sig []byte) bool {
	r, s, err := ParseDERSignature(sig)
	if err != nil {
		return false
	}
	return ecdsa.Verify(&x.PublicKey, hash, r, s)
}

func DERSignature(r, s *big.Int) []byte {
	buf := bytes.NewBuffer(nil)
	w := utils.NewBinaryWriter(buf)
	w.WriteByte(0x02)
	w.WriteDER(r.Bytes())
	w.WriteByte(0x02)
	w.WriteDER(s.Bytes())
	bin := buf.Bytes()
	buf = bytes.NewBuffer(nil)
	w = utils.NewBinaryWriter(buf)
	w.WriteByte(0x30)
	w.WriteDER(bin)
	return buf.Bytes()
}

func ParseDERSignature(sig []byte) (*big.Int, *big.Int, error) {
	buf := bytes.NewBuffer(sig)
	reader := utils.NewBinaryReader(buf)
	b, err := reader.ReadByte()
	if err != nil {
		return nil, nil, err
	}
	if b != 0x30 {
		return nil, nil, fmt.Errorf("invalid sig (got 0x%x instead of 0x30)", b)
	}
	pair, err := reader.ReadDER()
	if err != nil {
		return nil, nil, err
	}
	// prepare new
	buf = bytes.NewBuffer(pair)
	reader = utils.NewBinaryReader(buf)
	// R
	r, err := parseDERInt(reader)
	if err != nil {
		return nil, nil, err
	}
	// S
	s, err := parseDERInt(reader)
	if err != nil {
		return nil, nil, err
	}
	// finally
	return r, s, nil
}
func parseDERInt(reader utils.BinaryReader) (*big.Int, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != 0x02 {
		return nil, fmt.Errorf("invalid DER Int (got 0x%x instead of 0x02)", b)
	}
	raw, err := reader.ReadDER()
	if err != nil {
		return nil, err
	}
	return big.NewInt(0).SetBytes(raw), nil
}
