package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ensureDenomInPool check to make sure the input denoms exist in the provided pool asset map
func EnsureDenomInPool(poolAssetsByDenom map[string]PoolAsset, tokensIn sdk.Coins) error {
	for _, coin := range tokensIn {
		_, ok := poolAssetsByDenom[coin.Denom]
		if !ok {
			return errorsmod.Wrapf(ErrDenomNotFoundInPool, InvalidInputDenomsErrFormat, coin.Denom)
		}
	}

	return nil
}
