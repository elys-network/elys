package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName defines the module name
	ModuleName = "burner"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	LegacyHistoryKeyPrefix = "History/value/"
)

var (
	ParamsKeyPrefix  = []byte{0x01}
	HistoryKeyPrefix = []byte{0x02}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetHistoryKey(block uint64) []byte {
	return append(HistoryKeyPrefix, sdk.Uint64ToBigEndian(block)...)
}

func LegacyHistoryKey(
	timestamp string,
	denom string,
) []byte {
	var key []byte

	timestampBytes := []byte(timestamp)
	key = append(key, timestampBytes...)
	key = append(key, []byte("/")...)

	denomBytes := []byte(denom)
	key = append(key, denomBytes...)
	key = append(key, []byte("/")...)

	return key
}
