package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

type PoolAssetUSDValue struct {
	Asset string
	Value osmomath.BigDec
}

type InternalSwapRequest struct {
	InAmount sdk.Coin
	OutToken string
}

func (p *Pool) CalcJoinValueWithSlippage(ctx sdk.Context, oracleKeeper OracleKeeper,
	accountedPoolKeeper AccountedPoolKeeper, tokenIn sdk.Coin,
	weightMultiplier osmomath.BigDec, params Params) (osmomath.BigDec, osmomath.BigDec, error) {

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
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("token out denom not found")
	}

	outTokenPrice := oracleKeeper.GetDenomPrice(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("token price not set: %s", tokenOutDenom)
	}

	inTokenPrice := oracleKeeper.GetDenomPrice(ctx, tokenIn.Denom)
	if inTokenPrice.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("token price not set: %s", tokenIn.Denom)
	}

	joinValue := inTokenPrice.Mul(osmomath.BigDecFromSDKInt(tokenIn.Amount))

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOutDenom)
	if err != nil {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.LT(osmomath.OneBigDec()) {
		externalLiquidityRatio = osmomath.OneBigDec()
	}

	weightedAmount := osmomath.BigDecFromSDKInt(tokenIn.Amount).Mul(weightMultiplier)
	resizedAmount := osmomath.BigDecFromSDKInt(weightedAmount.Dec().TruncateInt()).
		Quo(externalLiquidityRatio).Dec().RoundInt()
	slippageAmount, err := p.CalcGivenInSlippage(
		ctx,
		oracleKeeper,
		p,
		sdk.Coins{sdk.NewCoin(tokenIn.Denom, resizedAmount)},
		tokenOutDenom,
		accountedPoolKeeper,
	)
	if err != nil {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	slippageAmount = slippageAmount.Mul(externalLiquidityRatio)
	slippageValue := slippageAmount.Mul(outTokenPrice)

	slippage := slippageValue.Quo(joinValue)

	minSlippage := params.GetBigDecMinSlippage().Mul(weightMultiplier)
	if slippage.LT(minSlippage) {
		slippage = minSlippage
		slippageValue = joinValue.Mul(minSlippage)
	}

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
	takerfees osmomath.BigDec,
) (tokensJoined sdk.Coins, numShares sdkmath.Int, slippage osmomath.BigDec, weightBalanceBonus osmomath.BigDec, swapFee osmomath.BigDec, takerFeesFinal osmomath.BigDec, err error) {
	// if it's not single sided liquidity, add at pool ratio
	if len(tokensIn) != 1 {
		numShares, tokensJoined, err := p.CalcJoinPoolNoSwapShares(tokensIn)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, osmomath.BigDec{}, osmomath.BigDec{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
		return tokensJoined, numShares, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil
	}

	if !p.PoolParams.UseOracle {
		tokenIn := tokensIn[0]
		totalSlippage := osmomath.ZeroBigDec()
		normalizedWeights := NormalizedWeights(p.PoolAssets)
		for _, weight := range normalizedWeights {
			if weight.Asset != tokenIn.Denom {
				_, slippage, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, weight.Asset, osmomath.ZeroBigDec(), accountedPoolKeeper)
				if err == nil {
					totalSlippage = totalSlippage.Add(slippage.Mul(weight.Weight))
				}
			}
		}

		numShares, tokensJoined, err := p.CalcSingleAssetJoinPoolShares(tokensIn)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, osmomath.BigDec{}, osmomath.BigDec{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
		return tokensJoined, numShares, totalSlippage, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil
	}

	accountedAssets := p.GetAccountedBalance(ctx, accountedPoolKeeper, p.PoolAssets)

	initialWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, accountedAssets, tokensIn[0].Denom)
	initialWeightOut := osmomath.OneBigDec().Sub(initialWeightIn)

	joinValueWithSlippage, slippage, err := p.CalcJoinValueWithSlippage(ctx, oracleKeeper, accountedPoolKeeper, tokensIn[0], initialWeightOut, params)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	tvl, err := p.TVL(ctx, oracleKeeper, accountedPoolKeeper)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	// Ensure tvl is not zero to avoid division by zero
	if tvl.IsZero() {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ErrAmountTooLow
	}

	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx, tokensIn, sdk.NewCoins(), accountedAssets)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	weightBalanceBonus, weightBreakingFee, isSwapFee := p.CalculateWeightFees(ctx, oracleKeeper, accountedAssets, newAssetPools, tokensIn[0].Denom, params, osmomath.OneBigDec())
	// apply percentage to fees, consider improvement or reduction of other token
	// Other denom weight ratio to reduce the weight breaking fees
	weightBreakingFee = weightBreakingFee.Mul(initialWeightOut)
	weightBalanceBonus = weightBalanceBonus.Mul(initialWeightOut)

	swapFee = osmomath.ZeroBigDec()
	if isSwapFee {
		swapFee = p.GetPoolParams().GetBigDecSwapFee().Mul(initialWeightOut)
	}

	takerFeesFinal = takerfees.Mul(initialWeightOut)

	totalShares := p.GetTotalShares()
	numSharesDec := osmomath.BigDecFromSDKInt(totalShares.Amount).
		Mul(joinValueWithSlippage).Quo(tvl).
		Mul(osmomath.OneBigDec().Sub(weightBreakingFee)).
		Mul(osmomath.OneBigDec().Sub(swapFee.Add(takerFeesFinal)))
	numShares = numSharesDec.Dec().RoundInt()
	err = p.IncreaseLiquidity(numShares, tokensIn)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	return tokensIn, numShares, slippage, weightBalanceBonus, swapFee, takerFeesFinal, nil
}
