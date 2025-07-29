package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/elys-network/elys/v7/utils"

	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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

// EnsureDenomInPool check to make sure the input denoms exist in the provided pool asset map
func EnsureDenomInPool(poolAssetsByDenom map[string]PoolAsset, tokensIn sdk.Coins) error {
	for _, coin := range tokensIn {
		_, ok := poolAssetsByDenom[coin.Denom]
		if !ok {
			return errorsmod.Wrapf(ErrDenomNotFoundInPool, InvalidInputDenomsErrFormat, coin.Denom)
		}
	}

	return nil
}

// ApplyDiscount applies discount to swap fee if applicable
func ApplyDiscount(swapFee osmomath.BigDec, discount osmomath.BigDec) osmomath.BigDec {
	// apply discount percentage to swap fee
	swapFee = swapFee.Mul(osmomath.OneBigDec().Sub(discount))
	return swapFee
}

// GetWeightBreakingFee When distanceDiff > 0, pool is getting worse so we calculate WBF based on final weight
// When distanceDiff < 0, pool is improving, we need to use initial weights. Say target is 50:50, initially pool is 80:20 and now after it is becoming 30:70,
// this is improving the pool but with finalWeightOut/finalWeightIn it will be 30/70 which doesn't provide enough bonus to user
func GetWeightBreakingFee(finalWeightIn, finalWeightOut, targetWeightIn, targetWeightOut, initialWeightIn, initialWeightOut osmomath.BigDec, distanceDiff osmomath.BigDec, params Params) osmomath.BigDec {
	weightBreakingFee := osmomath.ZeroBigDec()
	if !params.WeightBreakingFeeMultiplier.IsZero() {
		// (45/55*60/40) ^ 2.5
		if distanceDiff.IsPositive() {
			if !finalWeightOut.IsZero() && !finalWeightIn.IsZero() && !targetWeightOut.IsZero() && !targetWeightIn.IsZero() {
				weightBreakingFee = params.GetBigDecWeightBreakingFeeMultiplier().Mul(utils.Pow(finalWeightIn.Mul(targetWeightOut).Quo(finalWeightOut).Quo(targetWeightIn), params.GetBigDecWeightBreakingFeeExponent()))
			}
		} else {
			if !initialWeightOut.IsZero() && !initialWeightIn.IsZero() && !targetWeightOut.IsZero() && !targetWeightIn.IsZero() {
				weightBreakingFee = params.GetBigDecWeightBreakingFeeMultiplier().
					Mul(utils.Pow(initialWeightOut.Mul(targetWeightIn).Quo(initialWeightIn).Quo(targetWeightOut), params.GetBigDecWeightBreakingFeeExponent()))
			}
		}

		if weightBreakingFee.GT(osmomath.NewBigDecWithPrec(99, 2)) {
			weightBreakingFee = osmomath.NewBigDecWithPrec(99, 2)
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
		return errors.New("a token's weight in the pool must be greater than 0")
	}

	// TODO: add validation for asset weight overflow:
	// https://github.com/osmosis-labs/osmosis/issues/1958

	return nil
}
