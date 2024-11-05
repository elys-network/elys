package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Pool) CalcGivenOutSlippage(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokensOut sdk.Coins,
	tokenInDenom string,
	accPoolKeeper AccountedPoolKeeper,
) (sdk.Dec, error) {
	balancerInCoin, _, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, sdk.ZeroDec(), accPoolKeeper)
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
	ctx sdk.Context, oracleKeeper OracleKeeper, snapshot *Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec, accPoolKeeper AccountedPoolKeeper) (
	tokenIn sdk.Coin, slippage, slippageAmount sdk.Dec, weightBalanceBonus sdk.Dec, err error,
) {
	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		balancerInCoin, slippage, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, swapFee, accPoolKeeper)
		if err != nil {
			return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), err
		}

		err = p.applySwap(ctx, sdk.Coins{balancerInCoin}, tokensOut, swapFee, sdk.ZeroDec())
		if err != nil {
			return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), err
		}

		return balancerInCoin, slippage, sdk.ZeroDec(), sdk.ZeroDec(), nil
	}

	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOut.Denom)
	if outTokenPrice.IsZero() {
		return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	accountedAssets := p.GetAccountedBalance(ctx, accPoolKeeper, p.PoolAssets)
	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, accountedAssets)

	// in amount is calculated in this formula
	// balancer slippage amount = Max(oracleOutAmount-balancerOutAmount, 0)
	// resizedAmount = tokenIn / externalLiquidityRatio
	// actualSlippageAmount = balancer slippage(resizedAmount)
	oracleInAmount := sdk.NewDecFromInt(tokenOut.Amount).Mul(outTokenPrice).Quo(inTokenPrice)

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOut.Denom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), err
	}
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.IsZero() {
		return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), ErrAmountTooLow
	}

	resizedAmount := sdk.NewDecFromInt(tokenOut.Amount).Quo(externalLiquidityRatio).RoundInt()
	slippageAmount, err = p.CalcGivenOutSlippage(
		ctx,
		oracleKeeper,
		snapshot,
		sdk.Coins{sdk.NewCoin(tokenOut.Denom, resizedAmount)},
		tokenInDenom,
		accPoolKeeper,
	)
	if err != nil {
		return sdk.Coin{}, sdk.OneDec(), sdk.ZeroDec(), sdk.ZeroDec(), err
	}
	inAmountAfterSlippage := oracleInAmount.Add(slippageAmount)
	slippage = slippageAmount.Quo(inAmountAfterSlippage)

	// calculate weight distance difference to calculate bonus/cut on the operation
	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx,
		sdk.Coins{sdk.NewCoin(tokenInDenom, inAmountAfterSlippage.TruncateInt())},
		tokensOut, accountedAssets,
	)
	if err != nil {
		return sdk.Coin{}, sdk.OneDec(), sdk.ZeroDec(), sdk.ZeroDec(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// target weight
	// Asset weight remains same in new pool assets as in original pool assets
	targetWeightIn := GetDenomNormalizedWeight(newAssetPools, tokenInDenom)
	targetWeightOut := GetDenomNormalizedWeight(newAssetPools, tokenOut.Denom)

	// weight breaking fee as in Plasma pool
	weightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, newAssetPools, tokenInDenom)
	weightOut := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, newAssetPools, tokenOut.Denom)
	weightBreakingFee := GetWeightBreakingFee(weightIn, weightOut, targetWeightIn, targetWeightOut, p.PoolParams, distanceDiff)

	// weight recovery reward = weight breaking fee * weight recovery fee portion
	weightRecoveryReward := weightBreakingFee.Mul(p.PoolParams.WeightRecoveryFeePortion)

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus = weightBreakingFee.Neg()

	// If swap is improving weight, set weight breaking fee to zero
	if distanceDiff.IsNegative() {
		weightBreakingFee = sdk.ZeroDec()
		weightBalanceBonus = sdk.ZeroDec()

		// set weight breaking fee to zero if bonus is applied
		if initialWeightDistance.GT(p.PoolParams.ThresholdWeightDifference) {
			weightBalanceBonus = weightRecoveryReward
		}
	}

	if swapFee.GTE(sdk.OneDec()) {
		return sdk.Coin{}, sdk.OneDec(), sdk.ZeroDec(), sdk.ZeroDec(), ErrTooMuchSwapFee
	}

	tokenAmountInInt := inAmountAfterSlippage.
		Mul(sdk.OneDec().Add(weightBreakingFee)).
		Quo(sdk.OneDec().Sub(swapFee)).
		TruncateInt()
	tokenIn = sdk.NewCoin(tokenInDenom, tokenAmountInInt)
	err = p.applySwap(ctx, sdk.Coins{tokenIn}, tokensOut, swapFee, sdk.ZeroDec())
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), err
	}
	return tokenIn, slippage, slippageAmount, weightBalanceBonus, nil
}
