package types

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "perpetual"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_perpetual"

	// ParamsKey is the prefix for parameters of perpetual module
	ParamsKey = "perpetual_params"

	LegacyPoolKeyPrefix = "Pool/value/"
)

const MaxPageLimit = 100

const (
	InfinitePriceString    = "infinite"
	TakeProfitPriceDefault = "10000000000000000000000000000000000000000" // 10^40
)

var (
	MTPPrefix          = []byte{0x01}
	MTPCountPrefix     = []byte{0x02}
	OpenMTPCountPrefix = []byte{0x04}
	WhitelistPrefix    = []byte{0x05}
	PoolKeyPrefix      = []byte{0x06}
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

func GetMTPKey(addr sdk.AccAddress, id uint64) []byte {
	return append(MTPPrefix, append(address.MustLengthPrefix(addr), sdk.Uint64ToBigEndian(id)...)...)
}

func GetLegacyMTPKey(address string, id uint64) []byte {
	return append(MTPPrefix, append([]byte(address), sdk.Uint64ToBigEndian(id)...)...)
}

func GetPoolKey(index uint64) []byte {
	key := PoolKeyPrefix
	return append(key, sdk.Uint64ToBigEndian(index)...)
}

func legacyPoolKey(index uint64) []byte {
	var key []byte

	indexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(indexBytes, index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}

func GetLegacyPoolKey(index uint64) []byte {
	key := KeyPrefix(LegacyPoolKeyPrefix)
	return append(key, legacyPoolKey(index)...)
}

func GetMTPPrefixForAddress(address string) []byte {
	return append(MTPPrefix, []byte(address)...)
}
