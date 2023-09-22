package types

var ParamsKey = []byte{0x00}

const (
	ModuleName = "clock"

	StoreKey = ModuleName

	QuerierRoute = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_clock"
)
