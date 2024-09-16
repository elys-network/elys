package types

const (
	// ModuleName defines the module name
	ModuleName = "tradeshield"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tradeshield"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	ParamsKey                = []byte{0x01}
	PendingSpotOrderKey      = []byte{0x02}
	PendingSpotOrderCountKey = []byte{0x03}
)
