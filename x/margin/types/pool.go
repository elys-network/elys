package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewPool(poolId uint64) Pool {
	return Pool{
		AmmPoolId:                    poolId,
		Health:                       sdk.ZeroDec(),
		Enabled:                      true,
		Closed:                       false,
		ExternalLiabilities:          sdk.ZeroInt(),
		ExternalCustody:              sdk.ZeroInt(),
		NativeLiabilities:            sdk.ZeroInt(),
		NativeCustody:                sdk.ZeroInt(),
		InterestRate:                 sdk.ZeroDec(),
		NativeAssetBalance:           sdk.ZeroInt(),
		ExternalAssetBalance:         sdk.ZeroInt(),
		UnsettledExternalLiabilities: sdk.ZeroInt(),
		BlockInterestNative:          sdk.ZeroInt(),
	}
}
