package types

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// calcSingleAssetJoin - balance pAo.
func (p *Pool) calcSingleAssetJoin(tokenIn sdk.Coin, spreadFactor osmomath.BigDec, tokenInPoolAsset PoolAsset, totalShares sdkmath.Int) (numShares sdkmath.Int, err error) {
	totalWeight := p.TotalWeight
	if totalWeight.IsZero() {
		return sdkmath.ZeroInt(), errors.New("pool misconfigured, total weight = 0")
	}
	normalizedWeight := osmomath.BigDecFromSDKInt(tokenInPoolAsset.Weight).Quo(osmomath.BigDecFromSDKInt(totalWeight))
	poolShares, err := calcPoolSharesOutGivenSingleAssetIn(
		osmomath.BigDecFromSDKInt(tokenInPoolAsset.Token.Amount),
		normalizedWeight,
		osmomath.BigDecFromSDKInt(totalShares),
		osmomath.BigDecFromSDKInt(tokenIn.Amount),
		spreadFactor,
	)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	return poolShares.Dec().TruncateInt(), nil
}

// CalcSingleAssetJoinPoolShares calculates the number of shares created to join pool with the provided amount of `tokenIn`.
// The input tokens must either be:
// - a single token
// - contain exactly the same tokens as the pool contains
//
// It returns the number of shares created, the amount of coins actually joined into the pool
// (in case of not being able to fully join), or an error.
func (p *Pool) CalcSingleAssetJoinPoolShares(tokensIn sdk.Coins, takerFees osmomath.BigDec) (numShares sdkmath.Int, tokensJoined sdk.Coins, err error) {
	// get all 'pool assets' (aka current pool liquidity + balancer weight)
	poolAssetsByDenom, err := GetPoolAssetsByDenom(p.GetAllPoolAssets())
	if err != nil {
		return sdkmath.ZeroInt(), sdk.NewCoins(), err
	}

	err = EnsureDenomInPool(poolAssetsByDenom, tokensIn)
	if err != nil {
		return sdkmath.ZeroInt(), sdk.NewCoins(), err
	}

	// ensure that there aren't too many or too few assets in `tokensIn`
	if tokensIn.Len() != 1 {
		return sdkmath.ZeroInt(), sdk.NewCoins(), errors.New("pool only supports LP'ing with one asset")
	}

	// 2) Single token provided, so do single asset join and exit.
	totalShares := p.GetTotalShares()
	numShares, err = p.calcSingleAssetJoin(tokensIn[0], p.PoolParams.GetBigDecSwapFee().Add(takerFees), poolAssetsByDenom[tokensIn[0].Denom], totalShares.Amount)
	if err != nil {
		return sdkmath.ZeroInt(), sdk.NewCoins(), err
	}
	// we join all the tokens.
	tokensJoined = tokensIn
	return numShares, tokensJoined, nil
}
