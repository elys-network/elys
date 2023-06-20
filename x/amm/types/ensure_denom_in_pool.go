package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensureDenomInPool check to make sure the input denoms exist in the provided pool asset map
func EnsureDenomInPool(poolAssetsByDenom map[string]PoolAsset, tokensIn sdk.Coins) error {
	for _, coin := range tokensIn {
		_, ok := poolAssetsByDenom[coin.Denom]
		if !ok {
			return sdkerrors.Wrapf(ErrDenomNotFoundInPool, InvalidInputDenomsErrFormat, coin.Denom)
		}
	}

	return nil
}
