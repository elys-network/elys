package types

const (
	// ModuleName defines the module name
	ModuleName = "tokenomics"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	GenesisInflationKey = "GenesisInflation/value/"
)

var (
	ParamKeyPrefix = []byte{0x01}
)
