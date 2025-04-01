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
	AmmPoolKeyPrefix  = []byte{0x04}
	PoolPrefixKey     = []byte{0x05}
)

func GetDebtKey(owner sdk.AccAddress, poolId uint64) []byte {
	return append(DebtPrefixKey, append(address.MustLengthPrefix(owner), sdk.Uint64ToBigEndian(poolId)...)...)
}

func GetPoolKey(poolId uint64) []byte {
	return append(PoolPrefixKey, sdk.Uint64ToBigEndian(poolId)...)
}

func GetInterestKey(poolId uint64) []byte {
	return append(InterestPrefixKey, sdk.Uint64ToBigEndian(poolId)...)
}

func GetAmmPoolKey(id uint64) []byte {
	return append(AmmPoolKeyPrefix, sdk.Uint64ToBigEndian(id)...)
}
