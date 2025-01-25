package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
)

func (p Pool) CalcGivenOutSlippage(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokensOut sdk.Coins,
	tokenInDenom string,
	accPoolKeeper AccountedPoolKeeper,
) (elystypes.Dec34, error) {
	balancerInCoin, _, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, sdkmath.LegacyZeroDec(), accPoolKeeper)
	if err != nil {
		return elystypes.ZeroDec34(), err
	}

	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return elystypes.ZeroDec34(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return elystypes.ZeroDec34(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOut.Denom)
	if outTokenPrice.IsZero() {
		return elystypes.ZeroDec34(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	// in amount is calculated in this formula
	oracleInAmount := elystypes.NewDec34FromInt(tokenOut.Amount).
		Mul(outTokenPrice).
		Quo(inTokenPrice)
	balancerIn := elystypes.NewDec34FromInt(balancerInCoin.Amount)
	balancerSlippage := balancerIn.Sub(oracleInAmount)
	if balancerSlippage.IsNegative() {
		return elystypes.ZeroDec34(), nil
	}
	return balancerSlippage, nil
}

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
// weightBreakingFeePerpetualFactor should be 1 if perpetual is not the one calling this function
// Pool, and it's bank balances are updated in keeper.UpdatePoolForSwap
func (p *Pool) SwapInAmtGivenOut(
	ctx sdk.Context, oracleKeeper OracleKeeper, snapshot *Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdkmath.LegacyDec, accPoolKeeper AccountedPoolKeeper, weightBreakingFeePerpetualFactor sdkmath.LegacyDec, params Params) (
	tokenIn sdk.Coin, slippage, slippageAmount elystypes.Dec34, weightBalanceBonus elystypes.Dec34, oracleInAmount elystypes.Dec34, swapFeeFinal sdkmath.LegacyDec, err error,
) {
	// Fixed gas consumption per swap to prevent spam
	ctx.GasMeter().ConsumeGas(BalancerGasFeeForSwap, "balancer swap computation")

	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		balancerInCoin, slippage, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, swapFee, accPoolKeeper)
		if err != nil {
			return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
		}
		return balancerInCoin, slippage, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), swapFee, nil
	}

	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOut.Denom)
	if outTokenPrice.IsZero() {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	accountedAssets := p.GetAccountedBalance(ctx, accPoolKeeper, p.PoolAssets)

	// in amount is calculated in this formula
	// balancer slippage amount = Max(oracleOutAmount-balancerOutAmount, 0)
	// resizedAmount = tokenIn / externalLiquidityRatio
	// actualSlippageAmount = balancer slippage(resizedAmount)
	oracleInAmount = elystypes.NewDec34FromInt(tokenOut.Amount).
		Mul(outTokenPrice).
		Quo(inTokenPrice)

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOut.Denom)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
	}
	externalLiquidityRatioDec34 := elystypes.NewDec34FromLegacyDec(externalLiquidityRatio)
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatioDec34.IsZero() {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), ErrAmountTooLow
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
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
	}
	inAmountAfterSlippage := oracleInAmount.Add(slippageAmount.Mul(externalLiquidityRatioDec34))
	slippageAmount = slippageAmount.Mul(externalLiquidityRatioDec34)
	slippage = slippageAmount.Quo(oracleInAmount)

	// calculate weight distance difference to calculate bonus/cut on the operation
	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx,
		sdk.Coins{sdk.NewCoin(tokenInDenom, inAmountAfterSlippage.ToInt())},
		tokensOut, accountedAssets,
	)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
	}

	weightBalanceBonus, weightBreakingFee, isSwapFee := p.CalculateWeightFees(ctx, oracleKeeper, accountedAssets, newAssetPools, tokenInDenom, params, weightBreakingFeePerpetualFactor)
	if !isSwapFee {
		swapFee = sdkmath.LegacyZeroDec()
	}

	if swapFee.GTE(sdkmath.LegacyOneDec()) {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), ErrTooMuchSwapFee
	}

	swapFeeDec34 := elystypes.NewDec34FromLegacyDec(swapFee)

	// We deduct a swap fee on the input asset. The swap happens by following the invariant curve on the input * (1 - swap fee)
	// and then the swap fee is added to the pool.
	// Thus in order to give X amount out, we solve the invariant for the invariant input. However invariant input = (1 - swapfee) * trade input.
	// Therefore we divide by (1 - swapfee) here
	tokenAmountInInt := inAmountAfterSlippage.
		Quo(elystypes.OneDec34().Sub(weightBreakingFee)).
		Quo(elystypes.OneDec34().Sub(swapFeeDec34)).
		ToInt() // We round up tokenInAmt, as this is whats charged for the swap, for the precise amount out.
	// Otherwise, the pool would under-charge by this rounding error.
	tokenIn = sdk.NewCoin(tokenInDenom, tokenAmountInInt)
	return tokenIn, slippage, slippageAmount, weightBalanceBonus, oracleInAmount, swapFee, nil
}
