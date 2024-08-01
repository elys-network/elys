package types

const (
	// ModuleName defines the module name
	ModuleName = "stablestake"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_stablestake"
)

var DebtPrefixKey = []byte{0x01}
var InterestPrefixKey = []byte{0x02}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
