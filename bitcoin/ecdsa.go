package bitcoin

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"math/big"

	gcrypto "github.com/alivanz/go-crypto"
	"github.com/alivanz/go-utils"
	"github.com/btcsuite/btcd/btcec"
)

var (
	// Used in RFC6979 implementation when testing the nonce for correctness
	one = big.NewInt(1)

	// oneInitializer is used to fill a byte slice with byte 0x01.  It is provided
	// here to avoid the need to create it multiple times.
	oneInitializer = []byte{0x01}
)

type signer struct {
	*ecdsa.PrivateKey
}

func NewWallet(privkey []byte) (gcrypto.Wallet, error) {
	c := btcec.S256()
	D := big.NewInt(0)
	D.SetBytes(privkey)

	pk := new(ecdsa.PrivateKey)
	pk.PublicKey.Curve = btcec.S256()
	pk.D = D
	pk.PublicKey.X, pk.PublicKey.Y = c.ScalarBaseMult(D.Bytes())

	return &signer{pk}, nil
}

// hashToInt converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
// This is borrowed from crypto/ecdsa.
func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

// nonceRFC6979 generates an ECDSA nonce (`k`) deterministically according to RFC 6979.
// It takes a 32-byte hash as an input and returns 32-byte nonce to be used in ECDSA algorithm.
func nonceRFC6979(privkey *big.Int, hash []byte) *big.Int {

	curve := btcec.S256()
	q := curve.Params().N
	x := privkey
	alg := sha256.New

	qlen := q.BitLen()
	holen := alg().Size()
	rolen := (qlen + 7) >> 3
	bx := append(int2octets(x, rolen), bits2octets(hash, curve, rolen)...)

	// Step B
	v := bytes.Repeat(oneInitializer, holen)

	// Step C (Go zeroes the all allocated memory)
	k := make([]byte, holen)

	// Step D
	k = mac(alg, k, append(append(v, 0x00), bx...))

	// Step E
	v = mac(alg, k, v)

	// Step F
	k = mac(alg, k, append(append(v, 0x01), bx...))

	// Step G
	v = mac(alg, k, v)

	// Step H
	for {
		// Step H1
		var t []byte

		// Step H2
		for len(t)*8 < qlen {
			v = mac(alg, k, v)
			t = append(t, v...)
		}

		// Step H3
		secret := hashToInt(t, curve)
		if secret.Cmp(one) >= 0 && secret.Cmp(q) < 0 {
			return secret
		}
		k = mac(alg, k, append(v, 0x00))
		v = mac(alg, k, v)
	}
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func int2octets(v *big.Int, rolen int) []byte {
	out := v.Bytes()

	// left pad with zeros if it's too short
	if len(out) < rolen {
		out2 := make([]byte, rolen)
		copy(out2[rolen-len(out):], out)
		return out2
	}

	// drop most significant bytes if it's too long
	if len(out) > rolen {
		out2 := make([]byte, rolen)
		copy(out2, out[len(out)-rolen:])
		return out2
	}

	return out
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bits2octets(in []byte, curve elliptic.Curve, rolen int) []byte {
	z1 := hashToInt(in, curve)
	z2 := new(big.Int).Sub(z1, curve.Params().N)
	if z2.Sign() < 0 {
		return int2octets(z1, rolen)
	}
	return int2octets(z2, rolen)
}

// mac returns an HMAC of the given key and message.
func mac(alg func() hash.Hash, k, m []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(nil)
}

// generates a deterministic ECDSA signature according to RFC 6979 and BIP 62.
func (x *signer) Sign(hash []byte) (*big.Int, *big.Int, error) {
	// return ecdsa.Sign(rand.Reader, x.PrivateKey, hash)
	privkey := x.PrivateKey
	N := btcec.S256().N
	halfOrder := new(big.Int).Rsh(btcec.S256().N, 1)
	k := nonceRFC6979(privkey.D, hash)
	inv := new(big.Int).ModInverse(k, N)
	r, _ := privkey.Curve.ScalarBaseMult(k.Bytes())
	if r.Cmp(N) == 1 {
		r.Sub(r, N)
	}

	if r.Sign() == 0 {
		return nil, nil, errors.New("calculated R is zero")
	}

	e := hashToInt(hash, privkey.Curve)
	s := new(big.Int).Mul(privkey.D, r)
	s.Add(s, e)
	s.Mul(s, inv)
	s.Mod(s, N)

	if s.Cmp(halfOrder) == 1 {
		s.Sub(N, s)
	}
	if s.Sign() == 0 {
		return nil, nil, errors.New("calculated S is zero")
	}
	return r, s, nil
}

func (x *signer) PubKey() (ecdsa.PublicKey, error) {
	return x.PublicKey, nil
}

func (x *signer) Verify(hash []byte, r, s *big.Int) bool {
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
