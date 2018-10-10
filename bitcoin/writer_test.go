package bitcoin

import (
	"bytes"
	"encoding/hex"
	"log"
	"testing"
)

// https://klmoney.wordpress.com/bitcoin-dissecting-transactions-part-2-building-a-transaction-by-hand/
func TestWriteTxMsg(t *testing.T) {
	pubkeyhash, _ := hex.DecodeString("dd6cce9f255a8cc17bda8ba0373df8e861cb866e")
	unspents := []TxInExtended{NewTxIn("7e3ab0ea65b60f7d1ff4b231016fc958bc0766a46770410caa0a1855459b6e41", 0, 40000, "76a91499b1ebcfc11a13df5161aba8160460fe1601d54188ac")}
	changeaddr, _ := AddressParseBase58("1NAK3za9MkbAkkSBMLcvmhTD6etgB4Vhpr")
	// recipients := []TxOut{NewTxOut(P2PKH(dest.PubKeyHash()), 20000)}
	recipients := []TxOut{}
	txins, txouts, _ := NoTxFilter.TxFilter(P2PKH(changeaddr.PubKeyHash()), unspents, recipients, 20000)
	txwriter := NewTxWriterV1(pubkeyhash)
	bin, _ := txwriter.TxWrite(txins, txouts)
	{
		r := bytes.NewBuffer(bin)
		t.Log(hex.EncodeToString(r.Next(4)))
		t.Log(hex.EncodeToString(r.Next(1)))

		t.Log(hex.EncodeToString(r.Next(32)))
		t.Log(hex.EncodeToString(r.Next(4)))
		t.Log(hex.EncodeToString(r.Next(1)))
		t.Log(hex.EncodeToString(r.Next(25)))
		t.Log(hex.EncodeToString(r.Next(4)))

		t.Log(hex.EncodeToString(r.Next(1)))
		t.Log(hex.EncodeToString(r.Next(8)))
		t.Log(hex.EncodeToString(r.Next(1)))
		t.Log(hex.EncodeToString(r.Next(25)))

		t.Log(hex.EncodeToString(r.Next(4)))
		t.Log(hex.EncodeToString(r.Next(4)))
	}

	ref, _ := hex.DecodeString("0100000001416e9b4555180aaa0c417067a46607bc58c96f0131b2f41f7d0fb665eab03a7e000000001976a91499b1ebcfc11a13df5161aba8160460fe1601d54188acffffffff01204e0000000000001976a914e81d742e2c3c7acd4c29de090fc2c4d4120b2bf888ac0000000001000000")
	if !bytes.Equal(bin, ref) {
		t.Fail()
	}

	// Final TX Msg
	refhash, _ := hex.DecodeString("456f9e1b6184d770f1a240da9a3c4458e55b6b4ba2244dd21404db30b3131b94")
	hash := BitcoinHasher.Hash(bin)
	if !bytes.Equal(refhash, hash) {
		t.Fail()
	}

	privkey, _ := hex.DecodeString("3cd0560f5b27591916c643a0b7aa69d03839380a738d2e912990dcc573715d2c")
	wallet, err := NewWallet(privkey)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	pubkey, _ := wallet.PubKey()
	log.Print(hex.EncodeToString(pubkey.X.Bytes()))
	log.Print(hex.EncodeToString(pubkey.Y.Bytes()))
	r, s, _ := wallet.Sign(hash)
	if !wallet.Verify(hash, r, s) {
		t.Fail()
	}
}
