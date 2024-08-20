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

func GetPositionKey(creator sdk.AccAddress, id uint64) []byte {
	return append(PositionPrefix, append(address.MustLengthPrefix(creator), GetUint64Bytes(id)...)...)
}

func GetPositionPrefixForAddress(creator sdk.AccAddress) []byte {
	return append(PositionPrefix, address.MustLengthPrefix(creator)...)
}
