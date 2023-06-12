package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

func NewPoolAddress(poolId uint64) sdk.AccAddress {
	key := append([]byte("pool"), sdk.Uint64ToBigEndian(poolId)...)
	return address.Module(ModuleName, key)
}

func NewPoolRebalanceTreasury(poolId uint64) sdk.AccAddress {
	key := append([]byte("pool_treasury"), sdk.Uint64ToBigEndian(poolId)...)
	return address.Module(ModuleName, key)
}
