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

	ElysStakedKeyPrefix      = "ElysStaked/value/"
	ElysStakeChangeKeyPrefix = "ElysStakeChanged/value/"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// ElysStakedKey returns the store key to retrieve a ElysStaked from the address fields
func ElysStakedKey(address string) []byte {
	var key []byte

	addressBytes := []byte(address)
	key = append(key, addressBytes...)
	key = append(key, []byte("/")...)

	return key
}
