package bitcoin

import (
	"bytes"

	utils "github.com/alivanz/go-utils"
)

const (
	OP_DUP         byte = 0x76
	OP_HASH160     byte = 0xa9
	OP_CHECKSIG    byte = 0xac
	OP_EQUALVERIFY byte = 0x88
	OP_TRUE        byte = 0x51
)

// Standard Transaction to Bitcoin address (pay-to-pubkey-hash)
func P2PKH(pubkeyhash []byte) []byte {
	buf := bytes.NewBuffer([]byte{OP_DUP, OP_HASH160})
	x := utils.NewBinaryWriter(buf)
	x.WriteDER(pubkeyhash)
	x.Write([]byte{OP_EQUALVERIFY, OP_CHECKSIG})
	return buf.Bytes()
}

// Anyone-Can-Spend Outputs
func AnyoneCanSpent() []byte {
	return []byte{OP_TRUE}
}
