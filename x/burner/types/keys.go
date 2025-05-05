package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "burner"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

var (
	ParamsKeyPrefix  = []byte{0x01}
	HistoryKeyPrefix = []byte{0x02}
)

func GetHistoryKey(block uint64) []byte {
	return append(HistoryKeyPrefix, sdk.Uint64ToBigEndian(block)...)
}
