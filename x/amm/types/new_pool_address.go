package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

func NewPoolAddress(poolId uint64) sdk.AccAddress {
	poolIdModuleName := GetPoolIdModuleName(poolId)
	return address.Module(poolIdModuleName)
}

func GetPoolIdModuleName(poolId uint64) string {
	poolIdStr := strconv.FormatUint(poolId, 10)
	poolIdModuleName := ModuleName + "/pool/account/" + poolIdStr
	return poolIdModuleName
}

func NewPoolRebalanceTreasury(poolId uint64) sdk.AccAddress {
	key := append([]byte("pool_treasury"), sdk.Uint64ToBigEndian(poolId)...)
	return address.Module(ModuleName, key)
}

func NewPoolRevenueAddress(poolId uint64) sdk.AccAddress {
	key := append([]byte("pool_revenue"), sdk.Uint64ToBigEndian(poolId)...)
	return address.Module(ModuleName, key)
}
