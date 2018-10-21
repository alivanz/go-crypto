package bitcoin

type NetworkDesc struct {
	PubKeyHashCode byte
	ScriptHashCode byte
}

var (
	MainNetworkDesc    = NetworkDesc{0x00, 0x05}
	TestnetNetworkDesc = NetworkDesc{0x6f, 0xc4}
)
