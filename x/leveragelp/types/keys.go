package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "leveragelp"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
	// ParamsKey is the prefix for parameters of leveragelp module
	ParamsKey = "leveragelp_params"
)

const MaxPageLimit = 50000

var (
	LegacyPositionPrefix          = []byte{0x01}
	LegacyPositionCountPrefix     = []byte{0x02}
	LegacyOpenPositionCountPrefix = []byte{0x04}

	WhitelistPrefix       = []byte{0x05}
	PositionPrefix        = []byte{0x06}
	PositionCounterPrefix = []byte{0x07}

	ADLCounterKeyPrefix     = []byte{0x08}
	FallbackOffsetKeyPrefix = []byte{0x09}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetUint64Bytes returns the byte representation of the ID
func GetUint64Bytes(ID uint64) []byte {
	IDBz := make([]byte, 8)
	binary.BigEndian.PutUint64(IDBz, ID)
	return IDBz
}

// GetUint64FromBytes returns ID in uint64 format from a byte array
func GetUint64FromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

func GetWhitelistKey(acc sdk.AccAddress) []byte {
	return append(WhitelistPrefix, address.MustLengthPrefix(acc)...)
}

func GetPoolPrefixKey(poolId uint64) []byte {
	key := append(PositionPrefix, sdk.Uint64ToBigEndian(poolId)...)
	key = append(key, []byte("/")...)
	return key
}

func GetPoolCreatorPrefixKey(poolId uint64, creator sdk.AccAddress) []byte {
	key := GetPoolPrefixKey(poolId)
	key = append(key, address.MustLengthPrefix(creator)...)
	key = append(key, []byte("/")...)
	return key
}

func GetPositionKey(poolId uint64, creator sdk.AccAddress, id uint64) []byte {
	key := GetPoolCreatorPrefixKey(poolId, creator)
	key = append(key, sdk.Uint64ToBigEndian(id)...)
	return key
}

func GetPositionCounterKey(poolId uint64) []byte {
	return append(PositionCounterPrefix, sdk.Uint64ToBigEndian(poolId)...)
}

func GetADLCounterKey(poolId uint64) []byte {
	return append(ADLCounterKeyPrefix, sdk.Uint64ToBigEndian(poolId)...)
}

func GetLegacyPositionKey(creator sdk.AccAddress, id uint64) []byte {
	return append(LegacyPositionPrefix, append(address.MustLengthPrefix(creator), GetUint64Bytes(id)...)...)
}
