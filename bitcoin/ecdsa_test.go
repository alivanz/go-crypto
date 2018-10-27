package bitcoin

import (
	"encoding/hex"
	"testing"
)

func TestECDSA(t *testing.T) {
	// message hash
	hash, _ := hex.DecodeString("456f9e1b6184d770f1a240da9a3c4458e55b6b4ba2244dd21404db30b3131b94")
	// derive key
	privkey, _ := hex.DecodeString("3cd0560f5b27591916c643a0b7aa69d03839380a738d2e912990dcc573715d2c")
	wallet, err := NewWallet(privkey)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	// do signature
	r, s, _ := wallet.Sign(hash)
	sig := DERSignature(r, s)
	if !wallet.Verify(hash, r, s) {
		t.Log(hex.EncodeToString(sig))
		t.Fail()
	}
	// verify known valid signature
	sigref, _ := hex.DecodeString("304402201c3be71e1794621cbe3a7adec1af25f818f238f5796d47152137eba710f2174a02204f8fe667b696e30012ef4e56ac96afb830bddffee3b15d2e474066ab3aa39bad")
	r, s, _ = ParseDERSignature(sigref)
	if !wallet.Verify(hash, r, s) {
		t.Log("sigref")
		t.Log(hex.EncodeToString(privkey))
		pubkey, _ := wallet.PubKey()
		t.Log(hex.EncodeToString(pubkey.X.Bytes()))
		t.Log(hex.EncodeToString(pubkey.Y.Bytes()))
		t.Log(hex.EncodeToString(hash))
		t.Log(hex.EncodeToString(sigref))
		t.Fail()
	}
}

func TestECDSA2(t *testing.T) {
	// message hash
	hash, _ := hex.DecodeString("456f9e1b6184d770f1a240da9a3c4458e55b6b4ba2244dd21404db30b3131b94")
	// derive key
	privkey, _ := hex.DecodeString("3cd0560f5b27591916c643a0b7aa69d03839380a738d2e912990dcc573715d2c")
	wallet, err := NewWallet2(privkey)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	// do signature
	r, s, _ := wallet.Sign(hash)
	sig := DERSignature(r, s)
	if !wallet.Verify(hash, r, s) {
		t.Log(hex.EncodeToString(sig))
		t.Fail()
	}
	// verify known valid signature
	sigref, _ := hex.DecodeString("304402201c3be71e1794621cbe3a7adec1af25f818f238f5796d47152137eba710f2174a02204f8fe667b696e30012ef4e56ac96afb830bddffee3b15d2e474066ab3aa39bad")
	r, s, _ = ParseDERSignature(sigref)
	if !wallet.Verify(hash, r, s) {
		t.Log("sigref")
		t.Log(hex.EncodeToString(privkey))
		pubkey, _ := wallet.PubKey()
		t.Log(hex.EncodeToString(pubkey.X.Bytes()))
		t.Log(hex.EncodeToString(pubkey.Y.Bytes()))
		t.Log(hex.EncodeToString(hash))
		t.Log(hex.EncodeToString(sigref))
		t.Fail()
	}
}
