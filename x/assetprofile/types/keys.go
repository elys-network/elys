package types

const (
	// ModuleName defines the module name
	ModuleName = "assetprofile"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_assetprofile"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	ParamKeyPrefix = []byte{0x01}
)
