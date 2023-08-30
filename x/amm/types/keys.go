package types

const (
	// ModuleName defines the module name
	ModuleName = "amm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// TStoreKey defines the transient store key
	TStoreKey = "transient_amm"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
