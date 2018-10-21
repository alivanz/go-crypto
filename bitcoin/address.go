package bitcoin

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/alivanz/go-crypto/bitcoin/base58"
	"golang.org/x/crypto/ripemd160"
)

var (
	InvalidChecksum = fmt.Errorf("Invalid Checksum")
)

type Address [25]byte

// Parse base58 address
func AddressParseBase58(address string) (*Address, error) {
	var addr Address
	data, err := base58.Decode(address)
	if err != nil {
		return nil, err
	}
	if len(data) != 25 {
		return nil, fmt.Errorf("invalid data length")
	}
	checksum1 := data[21:25]
	checksum2 := BitcoinHasher.Hash(data[:21])[:4]
	if !bytes.Equal(checksum1, checksum2) {
		return nil, InvalidChecksum
	}
	copy(addr[:], data)
	return &addr, nil
}

// Turn big.Int into fixed length bytes
func BigIntToBytes(x *big.Int, length int) []byte {
	out := make([]byte, length)
	raw := x.Bytes()
	copy(out[length-len(raw):], raw)
	return out
}

// Serialize public key
func SerializePubkey(pubkey *ecdsa.PublicKey, compressed bool) []byte {
	buf := bytes.NewBuffer(nil)
	if compressed {
		if pubkey.Y.Bit(0) == 0 {
			// even
			buf.WriteByte(0x02)
		} else {
			// odd
			buf.WriteByte(0x03)
		}
		buf.Write(BigIntToBytes(pubkey.X, 32))
	} else {
		buf.WriteByte(0x04)
		buf.Write(BigIntToBytes(pubkey.X, 32))
		buf.Write(BigIntToBytes(pubkey.Y, 32))
	}
	return buf.Bytes()
}

// Hash SHA256 + RIPMD160
func Hash160(data []byte) []byte {
	// sha256
	hash1 := sha256.Sum256(data)
	// ripemd160
	ripmd := ripemd160.New()
	ripmd.Write(hash1[:])
	hash2 := ripmd.Sum(nil)
	return hash2[:]
}

// Turn serialized pubkey into bitcoin address
func PubKeyToAddress(addrtype byte, pubkey []byte) (*Address, error) {
	if len(pubkey) != 65 && len(pubkey) != 33 {
		return nil, fmt.Errorf("not valid pubkey length (65 or 33 bytes)")
	}
	return PubKeyHashToAddress(addrtype, Hash160(pubkey))
}

// Turn pubkeyhash into bitcoin address
func PubKeyHashToAddress(addrtype byte, pubkeyhash []byte) (*Address, error) {
	var addr Address
	if len(pubkeyhash) != 20 {
		return nil, fmt.Errorf("invalid pubkeyhash length")
	}
	addr[0] = addrtype
	copy(addr[1:21], pubkeyhash)
	checksum := BitcoinHasher.Hash(addr[:21])[:4]
	copy(addr[21:25], checksum)
	return &addr, nil
}

// Get address type, ie
func (addr *Address) AddrType() byte {
	return addr[0]
}

// Print base58 format
func (addr *Address) String() string {
	return base58.Encode(addr[:])
}
func (addr *Address) PubKeyHash() []byte {
	return addr[1:21]
}
