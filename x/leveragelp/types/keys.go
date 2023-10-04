package types

const (
	// ModuleName defines the module name
	ModuleName = "leveragelp"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_leveragelp"

	// ParamsKey is the prefix for parameters of leveragelp module
	ParamsKey = "leveragelp_params"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
