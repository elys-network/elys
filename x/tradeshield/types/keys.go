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

	// ParamsKey is the prefix for parameters of tradeshield module
	ParamsKey = "tradeshield_params"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
