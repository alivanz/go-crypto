package bitcoin

import (
	"encoding/hex"
	"testing"
)

func TestAddress(t *testing.T) {
	// Uncompressed pubkey
	if pubkeyx, err := hex.DecodeString("0450863AD64A87AE8A2FE83C1AF1A8403CB53F53E486D8511DAD8A04887E5B23522CD470243453A299FA9E77237716103ABC11A1DF38855ED6F2EE187E9C582BA6"); err != nil {
		t.Log(err)
		t.Fail()
		return
	} else if addr, err := PubKeyToAddress(MainNetworkDesc.PubKeyHashCode, pubkeyx); err != nil {
		t.Log(err)
		t.Fail()
		return
	} else if addr.String() != "16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM" {
		t.Logf("%v != 16UwLL9Risc3QfPqBUvKofHmBQ7wMtjvM", addr.String())
		t.Logf(hex.EncodeToString(addr.PubKeyHash()))
		t.Fail()
		return
	}

	// Compressed pubkey
	// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
	if pubkeyx, err := hex.DecodeString("0250863ad64a87ae8a2fe83c1af1a8403cb53f53e486d8511dad8a04887e5b2352"); err != nil {
		t.Log(err)
		t.Fail()
		return
	} else if addr, err := PubKeyToAddress(MainNetworkDesc.PubKeyHashCode, pubkeyx); err != nil {
		t.Log(err)
		t.Fail()
		return
	} else if addr.String() != "1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs" {
		t.Logf("%v != 1PMycacnJaSqwwJqjawXBErnLsZ7RkXUAs", addr.String())
		t.Logf(hex.EncodeToString(addr.PubKeyHash()))
		t.Fail()
		return
	}

	// Compressed pubkey
	// https://lists.linuxfoundation.org/pipermail/bitcoin-dev/2012-January/001039.html
	if pubkeyx, err := hex.DecodeString("034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa"); err != nil {
		t.Log(err)
		t.Fail()
		return
	} else if addr, err := PubKeyToAddress(MainNetworkDesc.PubKeyHashCode, pubkeyx); err != nil {
		t.Log(err)
		t.Fail()
		return
	} else if addr.String() != "1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9" {
		t.Logf("%v != 1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9", addr.String())
		t.Logf(hex.EncodeToString(addr.PubKeyHash()))
		t.Fail()
		return
	}

	// Parse address
	if addr, err := AddressParseBase58("1Q1pE5vPGEEMqRcVRMbtBK842Y6Pzo6nK9"); err != nil {
		t.Log(err)
		t.Fail()
	} else if addr2, err := PubKeyHashToAddress(addr.AddrType(), addr.PubKeyHash()); err != nil {
		t.Log(err)
		t.Fail()
	} else if addr.String() != addr2.String() {
		t.Log("address mismatch")
		t.Fail()
	}
}
