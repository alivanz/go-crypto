package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/alivanz/go-crypto/bitcoin/base58"
	"golang.org/x/crypto/ripemd160"
)

var (
	InvalidChecksum = fmt.Errorf("Invalid Checksum")
)

const (
	BitcoinPubKeyHash byte = 0x00
	BitcoinScriptHash byte = 0x05
)

type Address [25]byte

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
func PubKeyToAddress(addrtype byte, pubkey []byte) (*Address, error) {
	if len(pubkey) != 65 && len(pubkey) != 33 {
		return nil, fmt.Errorf("not valid pubkey length (65 or 33 bytes)")
	}
	// sha256
	hash1 := sha256.Sum256(pubkey)
	// ripemd160
	ripmd := ripemd160.New()
	ripmd.Write(hash1[:])
	hash2 := ripmd.Sum(nil)
	// copy
	return PubKeyHashToAddress(addrtype, hash2[:])
}

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
func (addr *Address) Version() byte {
	return addr[0]
}
func (addr *Address) String() string {
	return base58.Encode(addr[:])
}
func (addr *Address) PubKeyHash() []byte {
	return addr[1:21]
}
