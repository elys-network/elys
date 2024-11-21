package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetPoolShareDenom(poolId uint64) string {
	return fmt.Sprintf("amm/pool/%d", poolId)
}

func GetPoolIdFromShareDenom(shareDenom string) (uint64, error) {
	poolId, err := strconv.Atoi(strings.TrimPrefix(shareDenom, "amm/pool/"))
	if err != nil {
		return 0, err
	}
	return uint64(poolId), nil
}

// poolAssetsCoins returns all the coins corresponding to a slice of pool assets.
func poolAssetsCoins(assets []PoolAsset) sdk.Coins {
	coins := sdk.Coins{}
	for _, asset := range assets {
		coins = coins.Add(asset.Token)
	}
	return coins
}

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

// AbsDifferenceWithSign returns | a - b |, (a - b).sign()
// a is mutated and returned.
func AbsDifferenceWithSign(a, b sdkmath.LegacyDec) (sdkmath.LegacyDec, bool) {
	if a.GTE(b) {
		return a.SubMut(b), false
	} else {
		return a.NegMut().AddMut(b), true
	}
}

// ApplyDiscount applies discount to swap fee if applicable
func ApplyDiscount(swapFee sdkmath.LegacyDec, discount sdkmath.LegacyDec) sdkmath.LegacyDec {
	// apply discount percentage to swap fee
	swapFee = swapFee.Mul(sdkmath.LegacyOneDec().Sub(discount))
	return swapFee
}

func GetWeightBreakingFee(weightIn, weightOut, targetWeightIn, targetWeightOut sdkmath.LegacyDec, distanceDiff sdkmath.LegacyDec, params Params) sdkmath.LegacyDec {
	weightBreakingFee := sdkmath.LegacyZeroDec()
	if !weightOut.IsZero() && !weightIn.IsZero() && !targetWeightOut.IsZero() && !targetWeightIn.IsZero() && !params.WeightBreakingFeeMultiplier.IsZero() {
		// (45/55*60/40) ^ 2.5
		if distanceDiff.IsPositive() {
			weightBreakingFee = params.WeightBreakingFeeMultiplier.
				Mul(Pow(weightIn.Mul(targetWeightOut).Quo(weightOut).Quo(targetWeightIn), params.WeightBreakingFeeExponent))
		} else {
			weightBreakingFee = params.WeightBreakingFeeMultiplier.
				Mul(Pow(weightOut.Mul(targetWeightIn).Quo(weightIn).Quo(targetWeightOut), params.WeightBreakingFeeExponent))
		}

		if weightBreakingFee.GT(sdkmath.LegacyNewDecWithPrec(99, 2)) {
			weightBreakingFee = sdkmath.LegacyNewDecWithPrec(99, 2)
		}
	}
	return weightBreakingFee
}

// GetPoolAssetsByDenom return a mapping from pool asset
// denom to the pool asset itself. There must be no duplicates.
// Returns error, if any found.
func GetPoolAssetsByDenom(poolAssets []PoolAsset) (map[string]PoolAsset, error) {
	poolAssetsByDenom := make(map[string]PoolAsset)
	for _, poolAsset := range poolAssets {
		_, ok := poolAssetsByDenom[poolAsset.Token.Denom]
		if ok {
			return nil, fmt.Errorf(FormatRepeatingPoolAssetsNotAllowedErrFormat, poolAsset.Token.Denom)
		}

		poolAssetsByDenom[poolAsset.Token.Denom] = poolAsset
	}
	return poolAssetsByDenom, nil
}

func GetPoolAssetByDenom(assets []PoolAsset, denom string) (PoolAsset, bool) {
	for _, asset := range assets {
		if asset.Token.Denom == denom {
			return asset, true
		}
	}
	return PoolAsset{}, false
}

// validates a pool asset, to check if it has a valid weight.
func (pa PoolAsset) validateWeight() error {
	if pa.Weight.LTE(sdkmath.ZeroInt()) {
		return fmt.Errorf("a token's weight in the pool must be greater than 0")
	}

	// TODO: add validation for asset weight overflow:
	// https://github.com/osmosis-labs/osmosis/issues/1958

	return nil
}
