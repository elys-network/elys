package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PoolAssetUSDValue struct {
	Asset string
	Value sdkmath.LegacyDec
}

type InternalSwapRequest struct {
	InAmount sdk.Coin
	OutToken string
}

func (p *Pool) CalcJoinValueWithSlippage(ctx sdk.Context, oracleKeeper OracleKeeper,
	accountedPoolKeeper AccountedPoolKeeper, tokenIn sdk.Coin,
	weightMultiplier sdkmath.LegacyDec) (sdkmath.LegacyDec, sdkmath.LegacyDec, error) {

	// As this is 2 token pool, tokenOut will be
	tokenOutDenom := ""
	for _, asset := range p.PoolAssets {
		if asset.Token.Denom == tokenIn.Denom {
			continue
		}
		tokenOutDenom = asset.Token.Denom
	}
	// Not possible, but we might require this when we have pools with assets more than 2
	if tokenOutDenom == "" {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("token out denom not found")
	}

	joinValue := sdkmath.LegacyZeroDec()
	tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenIn.Denom)
	if tokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", tokenIn.Denom)
	}
	v := tokenPrice.Mul(sdkmath.LegacyNewDecFromInt(tokenIn.Amount))
	joinValue = joinValue.Add(v)

	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenIn.Denom)
	if inTokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", tokenIn.Denom)
	}

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOutDenom)
	if err != nil {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.LT(sdkmath.LegacyOneDec()) {
		externalLiquidityRatio = sdkmath.LegacyOneDec()
	}

	weightedAmount := sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Mul(weightMultiplier)
	resizedAmount := sdkmath.LegacyNewDecFromInt(weightedAmount.TruncateInt()).
		Quo(externalLiquidityRatio).RoundInt()
	slippageAmount, err := p.CalcGivenInSlippage(
		ctx,
		oracleKeeper,
		p,
		sdk.Coins{sdk.NewCoin(tokenIn.Denom, resizedAmount)},
		tokenOutDenom,
		accountedPoolKeeper,
	)
	if err != nil {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	slippageAmount = slippageAmount.Mul(externalLiquidityRatio)

	slippageValue := slippageAmount.Mul(inTokenPrice)
	slippage := slippageValue.Quo(joinValue)
	joinValueWithSlippage := joinValue.Sub(slippageValue)

	return joinValueWithSlippage, slippage, nil
}

// JoinPool calculates the number of shares needed for an all-asset join given tokensIn with swapFee applied.
// It updates the liquidity if the pool is joined successfully. If not, returns error.
func (p *Pool) JoinPool(
	ctx sdk.Context, snapshot *Pool,
	oracleKeeper OracleKeeper,
	accountedPoolKeeper AccountedPoolKeeper, tokensIn sdk.Coins,
	params Params,
) (tokensJoined sdk.Coins, numShares sdkmath.Int, slippage sdkmath.LegacyDec, weightBalanceBonus sdkmath.LegacyDec, swapFee sdkmath.LegacyDec, err error) {
	// if it's not single sided liquidity, add at pool ratio
	if len(tokensIn) != 1 {
		numShares, tokensJoined, err := p.CalcJoinPoolNoSwapShares(tokensIn)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, sdkmath.LegacyDec{}, sdkmath.LegacyDec{}, sdkmath.LegacyZeroDec(), err
		}
		return tokensJoined, numShares, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), nil
	}

	if !p.PoolParams.UseOracle {
		tokenIn := tokensIn[0]
		totalSlippage := sdkmath.LegacyZeroDec()
		normalizedWeights := NormalizedWeights(p.PoolAssets)
		for _, weight := range normalizedWeights {
			if weight.Asset != tokenIn.Denom {
				_, slippage, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, weight.Asset, sdkmath.LegacyZeroDec(), accountedPoolKeeper)
				if err == nil {
					totalSlippage = totalSlippage.Add(slippage.Mul(weight.Weight))
				}
			}
		}

		numShares, tokensJoined, err := p.CalcSingleAssetJoinPoolShares(tokensIn)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, sdkmath.LegacyDec{}, sdkmath.LegacyDec{}, sdkmath.LegacyZeroDec(), err
		}
		return tokensJoined, numShares, totalSlippage, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), nil
	}

	accountedAssets := p.GetAccountedBalance(ctx, accountedPoolKeeper, p.PoolAssets)

	initialWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, accountedAssets, tokensIn[0].Denom)
	initialWeightOut := sdkmath.LegacyOneDec().Sub(initialWeightIn)

	joinValueWithSlippage, slippage, err := p.CalcJoinValueWithSlippage(ctx, oracleKeeper, accountedPoolKeeper, tokensIn[0], initialWeightOut)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	tvl, err := p.TVL(ctx, oracleKeeper, accountedPoolKeeper)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// Ensure tvl is not zero to avoid division by zero
	if tvl.IsZero() {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrAmountTooLow
	}

	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx, tokensIn, sdk.NewCoins(), accountedAssets)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	weightBalanceBonus, weightBreakingFee, isSwapFee := p.CalculateWeightFees(ctx, oracleKeeper, accountedAssets, newAssetPools, tokensIn[0].Denom, params, sdkmath.LegacyOneDec())
	// apply percentage to fees, consider improvement or reduction of other token
	// Other denom weight ratio to reduce the weight breaking fees
	weightBreakingFee = weightBreakingFee.Mul(initialWeightOut)
	weightBalanceBonus = weightBalanceBonus.Mul(initialWeightOut)

	swapFee = sdkmath.LegacyZeroDec()
	if isSwapFee {
		swapFee = p.GetPoolParams().SwapFee.Mul(initialWeightOut)
	}

	totalShares := p.GetTotalShares()
	numSharesDec := sdkmath.LegacyNewDecFromInt(totalShares.Amount).
		Mul(joinValueWithSlippage).Quo(tvl).
		Mul(sdkmath.LegacyOneDec().Sub(weightBreakingFee)).
		Mul(sdkmath.LegacyOneDec().Sub(swapFee))
	numShares = numSharesDec.RoundInt()
	err = p.IncreaseLiquidity(numShares, tokensIn)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	return tokensIn, numShares, slippage, weightBalanceBonus, swapFee, nil
}
