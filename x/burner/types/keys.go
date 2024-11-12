package types

const (
	// ModuleName defines the module name
	ModuleName = "burner"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

var (
	ParamsKeyPrefix = []byte{0x01}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
