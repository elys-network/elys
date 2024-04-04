package types

const (
	// ModuleName defines the module name
	ModuleName = "estaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_estaking"

	// ParamsKey is the prefix to retrieve all Params
	ParamsKey = "Params/value/"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
