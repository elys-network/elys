package types

const (
	// ModuleName defines the module name
	ModuleName = "vaults"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_vaults"

    
)

var (
	ParamsKey = []byte("p_vaults")
)



func KeyPrefix(p string) []byte {
    return []byte(p)
}
