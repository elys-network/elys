package types

const (
	// ModuleName defines the module name
	ModuleName = "margin"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_margin"

	// ParamsKey is the prefix for parameters of margin module
	ParamsKey = "margin_params"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
