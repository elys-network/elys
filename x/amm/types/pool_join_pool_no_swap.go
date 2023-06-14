package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) TVL(ctx sdk.Context, oracleKeeper OracleKeeper) (sdk.Dec, error) {
	// TODO: handle non-oracle pool case
	// OracleAssetsTVL * TotalWeight / OracleAssetsWeight

	tvl := sdk.ZeroDec()
	for _, asset := range p.PoolAssets {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Token.Denom)
		if tokenPrice.IsZero() {
			return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", asset.Token.Denom)
		}
		v := tokenPrice.Mul(sdk.NewDecFromInt(asset.Token.Amount))
		tvl = tvl.Add(v)
	}
	return tvl, nil
}

// JoinPoolNoSwap calculates the number of shares needed for an all-asset join given tokensIn with swapFee applied.
// It updates the liquidity if the pool is joined successfully. If not, returns error.
func (p *Pool) JoinPoolNoSwap(ctx sdk.Context, oracleKeeper OracleKeeper, tokensIn sdk.Coins, swapFee sdk.Dec) (numShares math.Int, err error) {
	numShares, tokensJoined, err := p.CalcJoinPoolNoSwapShares(tokensIn, swapFee)
	if err != nil {
		return math.Int{}, err
	}

	if p.PoolParams.UseOracle {
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

		newAssetPools := p.NewPoolAssetsAfterSwap(
			tokensIn,
			sdk.Coins{},
		)
		weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)

		distanceDiff := weightDistance.Sub(initialWeightDistance)
		weightBreakingFee := sdk.ZeroDec()
		if distanceDiff.IsPositive() {
			weightBreakingFee = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff)
		}
		weightBalanceBonus := sdk.ZeroDec()
		if weightDistance.LT(p.PoolParams.ThresholdWeightDifference) && distanceDiff.IsNegative() {
			weightBalanceBonus = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff).Abs()
		}

		totalShares := p.GetTotalShares()
		numShares := sdk.NewDecFromInt(totalShares.Amount).
			Mul(joinValue).Quo(tvl).
			Mul(sdk.OneDec().Add(weightBalanceBonus).Sub(weightBreakingFee))
		return numShares.RoundInt(), nil
	}

	// update pool with the calculated share and liquidity needed to join pool
	p.IncreaseLiquidity(numShares, tokensJoined)
	return numShares, nil
}
