package types

import (
	"encoding/binary"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
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

const MaxPageLimit = 50000

var (
	PositionPrefix          = []byte{0x01}
	PositionCountPrefix     = []byte{0x02}
	OpenPositionCountPrefix = []byte{0x04}
	WhitelistPrefix         = []byte{0x05}
	SQBeginBlockPrefix      = []byte{0x06}
	LiquidationSortPrefix   = []byte{0x07} // Position liquidation sort prefix
	StopLossSortPrefix      = []byte{0x08} // Position stop loss sort prefix
	OffsetKeyPrefix         = []byte{0x09}
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

func GetLegacyWhitelistKey(address string) []byte {
	return append(WhitelistPrefix, []byte(address)...)
}

func GetPositionKey(creator sdk.AccAddress, id uint64) []byte {
	return append(PositionPrefix, append(address.MustLengthPrefix(creator), GetUint64Bytes(id)...)...)
}

func GetLegacyPositionKey(address string, id uint64) []byte {
	return append(PositionPrefix, append([]byte(address), GetUint64Bytes(id)...)...)
}

func GetPositionPrefixForAddress(creator sdk.AccAddress) []byte {
	return append(PositionPrefix, address.MustLengthPrefix(creator)...)
}

func GetLiquidationSortPrefix(poolId uint64) []byte {
	return append(LiquidationSortPrefix, GetUint64Bytes(poolId)...)
}

func GetLiquidationSortKey(poolId uint64, lpAmount math.Int, borrowed math.Int, id uint64) []byte {
	poolIdPrefix := GetLiquidationSortPrefix(poolId)
	if lpAmount.IsZero() || borrowed.IsZero() {
		return []byte{}
	}

	// default precision is 18
	// final string = decimalvalue + positionId(consistentlength)
	sortDec := math.LegacyNewDecFromInt(lpAmount).QuoInt(borrowed)
	paddedPosition := IntToStringWithPadding(id)
	bytes := []byte(sortDec.String() + paddedPosition)
	return append(poolIdPrefix, bytes...)
}

func IntToStringWithPadding(position uint64) string {
	// Define the desired length of the output string
	const length = 9

	// Convert the integer to a string
	str := strconv.FormatUint(position, 18)

	// Calculate the number of leading zeros needed
	padding := length - len(str)

	// Create the leading zeros string
	leadingZeros := ""
	for i := 0; i < padding; i++ {
		leadingZeros += "0"
	}

	// Concatenate leading zeros with the original number string
	result := leadingZeros + str
	return result
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
