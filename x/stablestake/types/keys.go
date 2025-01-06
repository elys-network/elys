package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName defines the module name
	ModuleName = "stablestake"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

var (
	DebtPrefixKey     = []byte{0x01}
	InterestPrefixKey = []byte{0x02}
	ParamKeyPrefix    = []byte{0x03}
	PoolPrefixKey     = []byte{0x04}
)

func GetDebtKey(owner sdk.AccAddress) []byte {
	return append(DebtPrefixKey, address.MustLengthPrefix(owner)...)
}

func GetPoolKey(poolId uint64) []byte {
	return append(PoolPrefixKey, sdk.Uint64ToBigEndian(poolId)...)
}

func GetInterestKey(poolId uint64) []byte {
	return append(InterestPrefixKey, sdk.Uint64ToBigEndian(poolId)...)
}

func GetDebtPrefixKey(poolId uint64) []byte {
	return append(DebtPrefixKey, sdk.Uint64ToBigEndian(poolId)...)
}
