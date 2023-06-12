package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (p *Pool) SwapInAmtGivenOut(
	ctx sdk.Context, oracleKeeper OracleKeeper, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (
	tokenIn sdk.Coin, weightBalanceBonus sdk.Dec, err error,
) {
	tokenInCoin, err := p.CalcInAmtGivenOut(tokensOut, tokenInDenom, swapFee)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		err = p.applySwap(ctx, sdk.Coins{tokenInCoin}, tokensOut)
		if err != nil {
			return sdk.Coin{}, sdk.ZeroDec(), err
		}
		return tokenInCoin, sdk.ZeroDec(), nil
	}

	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return sdk.Coin{}, sdk.ZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOut.Denom)
	if outTokenPrice.IsZero() {
		return sdk.Coin{}, sdk.ZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, p.PoolAssets)

	// in amount is calculated in this formula
	// slippage = (oracleOutAmount-balancerOutAmount)*slippageReduction
	// outAmountAfterSlippage = oracleOutAmount - slippage
	// TODO: consider when slippage is positive
	oracleInAmount := sdk.NewDecFromInt(tokenOut.Amount).Mul(outTokenPrice).Quo(inTokenPrice)
	balancerSlippage := oracleInAmount.Sub(sdk.NewDecFromInt(tokenInCoin.Amount))
	slippage := balancerSlippage.Mul(sdk.OneDec().Sub(p.PoolParams.SlippageReduction))
	inAmountAfterSlippage := oracleInAmount.Sub(slippage)

	// calculate weight distance difference to calculate bonus/cut on the operation
	newAssetPools := p.NewPoolAssetsAfterSwap(
		sdk.Coins{sdk.NewCoin(tokenInDenom, inAmountAfterSlippage.TruncateInt())},
		tokensOut,
	)
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// cut is valid when distance higher than original distance
	weightBreakingFee := sdk.ZeroDec()
	if distanceDiff.IsPositive() {
		weightBreakingFee = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff)
	}

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus = sdk.ZeroDec()
	// TODO: bonus should be coming from separate pool
	if weightDistance.LT(p.PoolParams.ThresholdWeightDifference) && distanceDiff.IsNegative() {
		weightBalanceBonus = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff).Abs()
		// TODO: we might skip swap fee in case it's a balance recovery operation
		// TODO: what if weightBalanceBonus amount is not enough since it's large swap? (Should provide maximum)
		// TODO: weightBalanceBonus should maintain several tokens - not just USD and swap out amount is in that token
	}
	tokenAmountInInt := inAmountAfterSlippage.Mul(sdk.OneDec().Add(weightBreakingFee)).TruncateInt()
	oracleInCoin := sdk.NewCoin(tokenInDenom, tokenAmountInInt)
	err = p.applySwap(ctx, sdk.Coins{oracleInCoin}, tokensOut)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}
	return oracleInCoin, weightBalanceBonus, nil
}
