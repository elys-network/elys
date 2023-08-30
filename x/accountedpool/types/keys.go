package types

const (
	// ModuleName defines the module name
	ModuleName = "tvl"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tvl"

	// ParamsKey is the prefix for parameters of margin module
	ParamsKey = "tvl_params"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
