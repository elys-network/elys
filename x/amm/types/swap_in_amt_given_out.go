package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Pool) CalcGivenOutSlippage(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokensOut sdk.Coins,
	tokenInDenom string,
	accPoolKeeper AccountedPoolKeeper,
) (sdkmath.LegacyDec, error) {
	balancerInCoin, _, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, sdkmath.LegacyZeroDec(), accPoolKeeper)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOut.Denom)
	if outTokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	// in amount is calculated in this formula
	oracleInAmount := sdkmath.LegacyNewDecFromInt(tokenOut.Amount).Mul(outTokenPrice).Quo(inTokenPrice)
	balancerIn := sdkmath.LegacyNewDecFromInt(balancerInCoin.Amount)
	balancerSlippage := balancerIn.Sub(oracleInAmount)
	if balancerSlippage.IsNegative() {
		return sdkmath.LegacyZeroDec(), nil
	}
	return balancerSlippage, nil
}

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
// weightBreakingFeePerpetualFactor should be 1 if perpetual is not the one calling this function
// Pool, and it's bank balances are updated in keeper.UpdatePoolForSwap
func (p *Pool) SwapInAmtGivenOut(
	ctx sdk.Context, oracleKeeper OracleKeeper, snapshot *Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdkmath.LegacyDec, accPoolKeeper AccountedPoolKeeper, weightBreakingFeePerpetualFactor sdkmath.LegacyDec, params Params, takerFees sdkmath.LegacyDec) (
	tokenIn sdk.Coin, slippage, slippageAmount sdkmath.LegacyDec, weightBalanceBonus sdkmath.LegacyDec, oracleInAmount sdkmath.LegacyDec, swapFeeFinal sdkmath.LegacyDec, err error,
) {
	// Fixed gas consumption per swap to prevent spam
	ctx.GasMeter().ConsumeGas(BalancerGasFeeForSwap, "balancer swap computation")

	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		balancerInCoin, slippage, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, swapFee, accPoolKeeper)
		if err != nil {
			return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}
		return balancerInCoin, slippage, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), swapFee, nil
	}

	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOut.Denom)
	if outTokenPrice.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	accountedAssets := p.GetAccountedBalance(ctx, accPoolKeeper, p.PoolAssets)

	// in amount is calculated in this formula
	// balancer slippage amount = Max(oracleOutAmount-balancerOutAmount, 0)
	// resizedAmount = tokenIn / externalLiquidityRatio
	// actualSlippageAmount = balancer slippage(resizedAmount)
	oracleInAmount = sdkmath.LegacyNewDecFromInt(tokenOut.Amount).Mul(outTokenPrice).Quo(inTokenPrice)

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOut.Denom)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrAmountTooLow
	}

	resizedAmount := sdkmath.LegacyNewDecFromInt(tokenOut.Amount).Quo(externalLiquidityRatio).RoundInt()
	slippageAmount, err = p.CalcGivenOutSlippage(
		ctx,
		oracleKeeper,
		snapshot,
		sdk.Coins{sdk.NewCoin(tokenOut.Denom, resizedAmount)},
		tokenInDenom,
		accPoolKeeper,
	)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	inAmountAfterSlippage := oracleInAmount.Add(slippageAmount.Mul(externalLiquidityRatio))
	slippageAmount = slippageAmount.Mul(externalLiquidityRatio)
	slippage = slippageAmount.Quo(oracleInAmount)

	if slippage.LT(params.MinSlippage) {
		slippage = params.MinSlippage
		slippageAmount = oracleInAmount.Mul(params.MinSlippage)
	}

	// calculate weight distance difference to calculate bonus/cut on the operation
	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx,
		sdk.Coins{sdk.NewCoin(tokenInDenom, inAmountAfterSlippage.TruncateInt())},
		tokensOut, accountedAssets,
	)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	weightBalanceBonus, weightBreakingFee, isSwapFee := p.CalculateWeightFees(ctx, oracleKeeper, accountedAssets, newAssetPools, tokenInDenom, params, weightBreakingFeePerpetualFactor)
	if !isSwapFee {
		swapFee = sdkmath.LegacyZeroDec()
	}

	if swapFee.GTE(sdkmath.LegacyOneDec()) {
		return sdk.Coin{}, sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrTooMuchSwapFee
	}

	// We deduct a swap fee on the input asset. The swap happens by following the invariant curve on the input * (1 - swap fee)
	// and then the swap fee is added to the pool.
	// Thus in order to give X amount out, we solve the invariant for the invariant input. However invariant input = (1 - swapfee) * trade input.
	// Therefore we divide by (1 - swapfee) here
	tokenAmountInInt := inAmountAfterSlippage.
		Quo(sdkmath.LegacyOneDec().Sub(weightBreakingFee)).
		Quo(sdkmath.LegacyOneDec().Sub(swapFee)).
		Ceil().TruncateInt() // We round up tokenInAmt, as this is whats charged for the swap, for the precise amount out.
	// Otherwise, the pool would under-charge by this rounding error.
	tokenIn = sdk.NewCoin(tokenInDenom, tokenAmountInInt)
	return tokenIn, slippage, slippageAmount, weightBalanceBonus, oracleInAmount, swapFee, nil
}
