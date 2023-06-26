package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// JoinPoolNoSwap calculates the number of shares needed for an all-asset join given tokensIn with swapFee applied.
// It updates the liquidity if the pool is joined successfully. If not, returns error.
func (p *Pool) JoinPoolNoSwap(ctx sdk.Context, oracleKeeper OracleKeeper, tokensIn sdk.Coins) (numShares math.Int, err error) {
	if !p.PoolParams.UseOracle {
		numShares, tokensJoined, err := p.CalcJoinPoolNoSwapShares(tokensIn)
		if err != nil {
			return math.Int{}, err
		}

		// update pool with the calculated share and liquidity needed to join pool
		p.IncreaseLiquidity(numShares, tokensJoined)
		return numShares, nil
	}

	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, p.PoolAssets)
	tvl, err := p.TVL(ctx, oracleKeeper)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	joinValue := sdk.ZeroDec()
	for _, asset := range tokensIn {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Denom)
		if tokenPrice.IsZero() {
			return sdk.ZeroInt(), fmt.Errorf("token price not set: %s", asset.Denom)
		}
		v := tokenPrice.Mul(sdk.NewDecFromInt(asset.Amount))
		joinValue = joinValue.Add(v)
	}

	newAssetPools, err := p.NewPoolAssetsAfterSwap(
		tokensIn,
		sdk.Coins{},
	)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)

	distanceDiff := weightDistance.Sub(initialWeightDistance)
	weightBreakingFee := sdk.ZeroDec()
	if distanceDiff.IsPositive() {
		weightBreakingFee = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff)
	}
	weightBalanceBonus := sdk.ZeroDec()
	if initialWeightDistance.GT(p.PoolParams.ThresholdWeightDifference) && distanceDiff.IsNegative() {
		weightBalanceBonus = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff).Abs()
	}

	totalShares := p.GetTotalShares()
	numSharesDec := sdk.NewDecFromInt(totalShares.Amount).
		Mul(joinValue).Quo(tvl).
		Mul(sdk.OneDec().Add(weightBalanceBonus).Sub(weightBreakingFee))
	numShares = numSharesDec.RoundInt()
	p.IncreaseLiquidity(numShares, tokensIn)
	return numShares, nil
}
