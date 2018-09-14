package bitcoin

import (
	"bytes"
	"fmt"

	"github.com/alivanz/go-crypto/bitcoin/base58"
)

var (
	InvalidChecksum = fmt.Errorf("Invalid Checksum")
)

type Address [25]byte

func (addr *Address) ParseBase58(address string) error {
	data, err := base58.Decode(address)
	if err != nil {
		return err
	}
	if len(data) != 25 {
		return fmt.Errorf("invalid data length")
	}
	checksum1 := data[21:25]
	checksum2 := BitcoinHasher.Hash(data[:21])[:4]
	if !bytes.Equal(checksum1, checksum2) {
		return InvalidChecksum
	}
	copy(addr[:], data)
	return nil
}
func (addr *Address) ParsePubKeyHash(pubkeyhash []byte) error {
	if len(pubkeyhash) != 20 {
		return fmt.Errorf("invalid pubkeyhash length")
	}
	copy(addr[1:21], pubkeyhash)
	checksum := BitcoinHasher.Hash(addr[:21])[:4]
	copy(addr[21:25], checksum)
	return nil
}
func (addr *Address) Network() byte {
	return addr[0]
}
func (addr *Address) String() string {
	return base58.Encode(addr[:])
}
func (addr *Address) PubKeyHash() []byte {
	return addr[1:21]
}
