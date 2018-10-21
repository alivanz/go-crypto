package bitcoin

import (
	"bytes"
	"math/big"

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

// P2PKH script signature
func P2PKHScriptSig(r, s *big.Int, pubkey []byte) []byte {
	buf := bytes.NewBuffer(nil)
	w := utils.NewBinaryWriter(buf)
	w.WriteCompact(DERSignatureFull(r, s))
	w.WriteCompact(pubkey)
	return buf.Bytes()
}

// DER Signature with sighash 0x01
func DERSignatureFull(r, s *big.Int) []byte {
	buf := bytes.NewBuffer(nil)
	w := utils.NewBinaryWriter(buf)
	w.WriteByte(0x30)
	w.WriteCompact(DERSignatureHalf(r, s))
	// Signature hash code
	w.WriteByte(0x01)
	return buf.Bytes()
}

// DER Signature, thats all
func DERSignatureHalf(r, s *big.Int) []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(DERInt(r))
	buf.Write(DERInt(s))
	return buf.Bytes()
}

// DER format for int with prefix 0x02
func DERInt(x *big.Int) []byte {
	buf := bytes.NewBuffer(nil)
	w := utils.NewBinaryWriter(buf)
	w.WriteByte(0x02)
	if x.BitLen()%8 == 0 {
		w.WriteCompact(append([]byte{0}, x.Bytes()...))
	} else {
		w.WriteCompact(x.Bytes())
	}
	return buf.Bytes()
}
