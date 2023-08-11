package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Pool) CalcGivenOutSlippage(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	tokensOut sdk.Coins,
	tokenInDenom string,
) (sdk.Dec, error) {
	balancerInCoin, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, tokensOut, tokenInDenom, sdk.ZeroDec())
	if err != nil {
		return sdk.ZeroDec(), err
	}

	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return sdk.ZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOut.Denom)
	if outTokenPrice.IsZero() {
		return sdk.ZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	// in amount is calculated in this formula
	oracleInAmount := sdk.NewDecFromInt(tokenOut.Amount).Mul(outTokenPrice).Quo(inTokenPrice)
	balancerIn := sdk.NewDecFromInt(balancerInCoin.Amount)
	balancerSlippage := balancerIn.Sub(oracleInAmount)
	if balancerSlippage.IsNegative() {
		return sdk.ZeroDec(), nil
	}
	return balancerSlippage, nil
}

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (p *Pool) SwapInAmtGivenOut(
	ctx sdk.Context, oracleKeeper OracleKeeper, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (
	tokenIn sdk.Coin, weightBalanceBonus sdk.Dec, err error,
) {
	balancerInCoin, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, tokensOut, tokenInDenom, swapFee)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		err = p.applySwap(ctx, sdk.Coins{balancerInCoin}, tokensOut, swapFee, sdk.ZeroDec())
		if err != nil {
			return sdk.Coin{}, sdk.ZeroDec(), err
		}
		return balancerInCoin, sdk.ZeroDec(), nil
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
	// balancer slippage amount = Max(oracleOutAmount-balancerOutAmount, 0)
	// resizedAmount = tokenIn / externalLiquidityRatio
	// actualSlippageAmount = balancer slippage(resizedAmount)
	oracleInAmount := sdk.NewDecFromInt(tokenOut.Amount).Mul(outTokenPrice).Quo(inTokenPrice)

	resizedAmount := sdk.NewDecFromInt(tokenOut.Amount).Quo(p.PoolParams.ExternalLiquidityRatio).RoundInt()
	slippageAmount, err := p.CalcGivenOutSlippage(ctx, oracleKeeper, sdk.Coins{sdk.NewCoin(tokenOut.Denom, resizedAmount)}, tokenInDenom)
	inAmountAfterSlippage := oracleInAmount.Add(slippageAmount)

	// calculate weight distance difference to calculate bonus/cut on the operation
	newAssetPools, err := p.NewPoolAssetsAfterSwap(
		sdk.Coins{sdk.NewCoin(tokenInDenom, inAmountAfterSlippage.TruncateInt())},
		tokensOut,
	)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// cut is valid when distance higher than original distance
	weightBreakingFee := sdk.ZeroDec()
	if distanceDiff.IsPositive() {
		weightBreakingFee = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff)
	}

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus = sdk.ZeroDec()
	if initialWeightDistance.GT(p.PoolParams.ThresholdWeightDifference) && distanceDiff.IsNegative() {
		weightBalanceBonus = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff).Abs()
		// TODO: we might skip swap fee in case it's a balance recovery operation
	}
	tokenAmountInInt := inAmountAfterSlippage.
		Mul(sdk.OneDec().Add(weightBreakingFee)).
		Quo(sdk.OneDec().Sub(swapFee)).
		TruncateInt()
	oracleInCoin := sdk.NewCoin(tokenInDenom, tokenAmountInInt)
	err = p.applySwap(ctx, sdk.Coins{oracleInCoin}, tokensOut, swapFee, sdk.ZeroDec())
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}
	return oracleInCoin, weightBalanceBonus, nil
}
