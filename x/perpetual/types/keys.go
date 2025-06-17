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
	MTPPrefix              = []byte{0x01}
	PerpetualCounterPrefix = []byte{0x02}
	WhitelistPrefix        = []byte{0x05}
	PoolKeyPrefix          = []byte{0x06}
	InterestRatePrefix     = []byte{0x07}
	FundingRatePrefix      = []byte{0x08}

	// Delete after v18 migration
	LegacyMTPCountPrefix     = []byte{0x02}
	LegacyOpenMTPCountPrefix = []byte{0x04}
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

func GetWhitelistKey(addr sdk.AccAddress) []byte {
	return append(WhitelistPrefix, address.MustLengthPrefix(addr)...)
}

func GetMTPPrefixForAddress(addr sdk.AccAddress) []byte {
	return append(MTPPrefix, address.MustLengthPrefix(addr)...)
}

func GetMTPPrefixForAddressAndPoolId(addr sdk.AccAddress, poolId uint64) []byte {
	key := GetMTPPrefixForAddress(addr)
	key = append(key, []byte("/")...)
	key = append(key, GetUint64Bytes(poolId)...)
	return key
}

func GetMTPKey(addr sdk.AccAddress, poolId, id uint64) []byte {
	key := GetMTPPrefixForAddressAndPoolId(addr, poolId)
	key = append(key, []byte("/")...)
	key = append(key, GetUint64Bytes(id)...)
	return key
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

func GePerpetualCounterKey(ammPoolId uint64) []byte {
	return append(PerpetualCounterPrefix, sdk.Uint64ToBigEndian(ammPoolId)...)
}
