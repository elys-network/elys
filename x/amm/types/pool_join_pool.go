package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
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
	weightMultiplier elystypes.Dec34, params Params) (elystypes.Dec34, elystypes.Dec34, error) {

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
		return elystypes.ZeroDec34(), elystypes.ZeroDec34(), fmt.Errorf("token out denom not found")
	}

	outTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return elystypes.ZeroDec34(), elystypes.ZeroDec34(), fmt.Errorf("token price not set: %s", tokenOutDenom)
	}

	inTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenIn.Denom)
	if inTokenPrice.IsZero() {
		return elystypes.ZeroDec34(), elystypes.ZeroDec34(), fmt.Errorf("token price not set: %s", tokenIn.Denom)
	}

	joinValue := inTokenPrice.MulInt(tokenIn.Amount)

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOutDenom)
	if err != nil {
		return elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.LT(sdkmath.LegacyOneDec()) {
		externalLiquidityRatio = sdkmath.LegacyOneDec()
	}

	weightedAmount := weightMultiplier.MulInt(tokenIn.Amount)
	resizedAmount := sdkmath.LegacyNewDecFromInt(weightedAmount.ToInt()).
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
		return elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}
	slippageAmount = slippageAmount.Mul(elystypes.NewDec34FromLegacyDec(externalLiquidityRatio))
	slippageValue := slippageAmount.Mul(outTokenPrice)

	slippage := slippageValue.Quo(joinValue)

	minSlippage := elystypes.NewDec34FromLegacyDec(params.MinSlippage).Mul(weightMultiplier)
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
	takerfees elystypes.Dec34,
) (tokensJoined sdk.Coins, numShares sdkmath.Int, slippage elystypes.Dec34, weightBalanceBonus elystypes.Dec34, swapFee elystypes.Dec34, takerFeesFinal elystypes.Dec34, err error) {
	// if it's not single sided liquidity, add at pool ratio
	if len(tokensIn) != 1 {
		numShares, tokensJoined, err := p.CalcJoinPoolNoSwapShares(tokensIn)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, elystypes.Dec34{}, elystypes.Dec34{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}
		return tokensJoined, numShares, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), nil
	}

	if !p.PoolParams.UseOracle {
		tokenIn := tokensIn[0]
		totalSlippage := elystypes.ZeroDec34()
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
			return sdk.NewCoins(), sdkmath.Int{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, elystypes.Dec34{}, elystypes.Dec34{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}
		return tokensJoined, numShares, totalSlippage, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), nil
	}

	accountedAssets := p.GetAccountedBalance(ctx, accountedPoolKeeper, p.PoolAssets)

	initialWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, accountedAssets, tokensIn[0].Denom)
	initialWeightOut := elystypes.OneDec34().Sub(initialWeightIn)

	joinValueWithSlippage, slippage, err := p.CalcJoinValueWithSlippage(ctx, oracleKeeper, accountedPoolKeeper, tokensIn[0], initialWeightOut, params)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}

	tvl, err := p.TVL(ctx, oracleKeeper, accountedPoolKeeper)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}

	// Ensure tvl is not zero to avoid division by zero
	if tvl.IsZero() {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), ErrAmountTooLow
	}

	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx, tokensIn, sdk.NewCoins(), accountedAssets)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}

	weightBalanceBonus, weightBreakingFee, isSwapFee := p.CalculateWeightFees(ctx, oracleKeeper, accountedAssets, newAssetPools, tokensIn[0].Denom, params, sdkmath.LegacyOneDec())
	// apply percentage to fees, consider improvement or reduction of other token
	// Other denom weight ratio to reduce the weight breaking fees
	weightBreakingFee = weightBreakingFee.Mul(initialWeightOut)
	weightBalanceBonus = weightBalanceBonus.Mul(initialWeightOut)

	swapFee = elystypes.ZeroDec34()
	if isSwapFee {
		swapFee = elystypes.NewDec34FromLegacyDec(p.GetPoolParams().SwapFee).Mul(initialWeightOut)
	}

	takerFeesFinal = takerfees.Mul(initialWeightOut)

	totalShares := p.GetTotalShares()
	numSharesDec := elystypes.NewDec34FromInt(totalShares.Amount).
		Mul(joinValueWithSlippage).Quo(tvl).
		Mul(elystypes.OneDec34().Sub(weightBreakingFee)).
		Mul(elystypes.OneDec34().Sub(swapFee.Add(takerFeesFinal)))
	numShares = numSharesDec.ToInt()
	err = p.IncreaseLiquidity(numShares, tokensIn)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}

	return tokensIn, numShares, slippage, weightBalanceBonus, swapFee, takerFeesFinal, nil
}
