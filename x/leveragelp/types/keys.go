package types

import (
	"encoding/binary"

	"cosmossdk.io/math"
)

const (
	// ModuleName defines the module name
	ModuleName = "leveragelp"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_leveragelp"

	// ParamsKey is the prefix for parameters of leveragelp module
	ParamsKey = "leveragelp_params"
)

const MaxPageLimit = 100

var (
	PositionPrefix          = []byte{0x01}
	PositionCountPrefix     = []byte{0x02}
	OpenPositionCountPrefix = []byte{0x04}
	WhitelistPrefix         = []byte{0x05}
	SQBeginBlockPrefix      = []byte{0x06}
	LiquidationSortPrefix   = []byte{0x07} // Position liquidation sort prefix
	StopLossSortPrefix      = []byte{0x08} // Position stop loss sort prefix
	PositionSortPrefix      = []byte{0x08} // Position sort prefix
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

func GetWhitelistKey(address string) []byte {
	return append(WhitelistPrefix, []byte(address)...)
}

func GetPositionKey(address string, id uint64) []byte {
	return append(PositionPrefix, append([]byte(address), GetUint64Bytes(id)...)...)
}

func GetLiquidationSortPrefix(poolId uint64) []byte {
	return append(LiquidationSortPrefix, GetUint64Bytes(poolId)...)
}

func GetLiquidationSortKey(poolId uint64, lpAmount math.Int, borrowed math.Int, id uint64) []byte {
	poolIdPrefix := GetLiquidationSortPrefix(poolId)
	if lpAmount.IsZero() || borrowed.IsZero() {
		return []byte{}
	}

	sortDec := math.LegacyNewDecFromInt(lpAmount).QuoInt(borrowed)
	bytes := sortDec.BigInt().Bytes()
	lengthPrefix := GetUint64Bytes(uint64(len(bytes)))
	posIdSuffix := GetUint64Bytes(id)
	return append(append(append(poolIdPrefix, lengthPrefix...), bytes...), posIdSuffix...)
}

func GetStopLossSortPrefix(poolId uint64) []byte {
	return append(StopLossSortPrefix, GetUint64Bytes(poolId)...)
}

func GetStopLossSortKey(poolId uint64, stopLossPrice math.LegacyDec, id uint64) []byte {
	poolIdPrefix := GetStopLossSortPrefix(poolId)
	if stopLossPrice.IsNil() || !stopLossPrice.IsPositive() {
		return []byte{}
	}

	bytes := stopLossPrice.BigInt().Bytes()
	lengthPrefix := GetUint64Bytes(uint64(len(bytes)))
	posIdSuffix := GetUint64Bytes(id)
	return append(append(append(poolIdPrefix, lengthPrefix...), bytes...), posIdSuffix...)
}

func GetPositionPrefixForAddress(address string) []byte {
	return append(PositionPrefix, []byte(address)...)
}
