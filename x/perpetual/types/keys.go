package types

import (
	"encoding/binary"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "perpetual"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// ParamsKey is the prefix for parameters of perpetual module
	ParamsKey = "perpetual_params"

	LegacyPoolKeyPrefix = "Pool/value/"
)

const MaxPageLimit = 10000

const (
	InfinitePriceString = "infinite"
	ZeroPriceString     = "zero"
)

var (
	TakeProfitPriceDefault = math.LegacyMustNewDecFromStr("10000000000000000000000000000000000000000") // 10^40
	StopLossPriceDefault   = math.LegacyZeroDec()
)

var (
	MTPPrefix          = []byte{0x01}
	MTPCountPrefix     = []byte{0x02}
	OpenMTPCountPrefix = []byte{0x04}
	WhitelistPrefix    = []byte{0x05}
	PoolKeyPrefix      = []byte{0x06}

	InterestRatePrefix = []byte{0x07}
	FundingRatePrefix  = []byte{0x08}
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

func GetWhitelistKey(addr sdk.AccAddress) []byte {
	return append(WhitelistPrefix, address.MustLengthPrefix(addr)...)
}

func GetLegacyWhitelistKey(address string) []byte {
	return append(WhitelistPrefix, []byte(address)...)
}

func GetMTPKey(addr sdk.AccAddress, id uint64) []byte {
	return append(MTPPrefix, append(address.MustLengthPrefix(addr), sdk.Uint64ToBigEndian(id)...)...)
}

func GetMTPPrefixForAddress(addr sdk.AccAddress) []byte {
	return append(MTPPrefix, address.MustLengthPrefix(addr)...)
}

func GetPoolKey(index uint64) []byte {
	key := PoolKeyPrefix
	return append(key, sdk.Uint64ToBigEndian(index)...)
}

func GetInterestRateKey(block uint64, pool uint64) []byte {
	return append(GetUint64Bytes(block), GetUint64Bytes(pool)...)
}

func GetFundingRateKey(block uint64, pool uint64) []byte {
	return append(GetUint64Bytes(block), GetUint64Bytes(pool)...)
}
