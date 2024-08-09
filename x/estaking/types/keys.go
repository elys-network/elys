package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "estaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_estaking"

	LegacyParamsKey                = "Params/value/"
	LegacyElysStakedKeyPrefix      = "ElysStaked/value/"
	LegacyElysStakeChangeKeyPrefix = "ElysStakeChanged/value/"
)

var (
	ElysStakedKeyPrefix      = []byte{0x01}
	ElysStakeChangeKeyPrefix = []byte{0x02}
	ParamsKeyPrefix          = []byte{0x03}
)

// remove after migration
func LegacyKeyPrefix(p string) []byte {
	return []byte(p)
}

func GetElysStakedKey(acc sdk.AccAddress) []byte {
	return append(ElysStakedKeyPrefix, address.MustLengthPrefix(acc)...)
}

func GetElysStakeChangeKey(acc sdk.AccAddress) []byte {
	return append(ElysStakeChangeKeyPrefix, address.MustLengthPrefix(acc)...)
}

// remove after migration
func LegacyElysStakedKey(address string) []byte {
	var key []byte

	addressBytes := []byte(address)
	key = append(key, addressBytes...)
	key = append(key, []byte("/")...)

	return key
}
