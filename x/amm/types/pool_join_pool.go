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

func (p *Pool) CalcJoinValueWithSlippage(ctx sdk.Context, snapshot SnapshotPool, oracleKeeper OracleKeeper, tokenIn sdk.Coin,
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
	resizedAmount := weightedAmount.
		Quo(externalLiquidityRatio).Dec().RoundInt()
	slippageAmount, _, err := p.CalcGivenInSlippage(
		ctx,
		oracleKeeper,
		snapshot,
		sdk.Coins{sdk.NewCoin(tokenIn.Denom, resizedAmount)},
		tokenOutDenom,
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
	ctx sdk.Context, snapshot SnapshotPool,
	oracleKeeper OracleKeeper,
	accountedPoolKeeper AccountedPoolKeeper, tokensIn sdk.Coins,
	params Params,
	takerFees osmomath.BigDec,
) (tokensJoined sdk.Coins, numShares sdkmath.Int, slippage osmomath.BigDec, weightBalanceBonus osmomath.BigDec, swapFee osmomath.BigDec, takerFeesFinal osmomath.BigDec, swapsInfos []SwapInfo, err error) {
	// if it's not single sided liquidity, add at pool ratio
	if len(tokensIn) != 1 {
		// We calculate based on snapshot, if there no accounted pool then it will be same as normal pool
		numShares, tokensJoined, err = snapshot.CalcJoinPoolNoSwapShares(tokensIn)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
		}

		// update pool with the calculated share and liquidity needed to join pool
		// if it's accounted pool, we increase liquidity on the original pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, osmomath.BigDec{}, osmomath.BigDec{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
		}
		return tokensJoined, numShares, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), swapsInfos, nil
	}

	if !p.PoolParams.UseOracle {
		tokenIn := tokensIn[0]
		totalSlippage := osmomath.ZeroBigDec()
		normalizedWeights := NormalizedWeights(p.PoolAssets)
		tokenAmountOut := osmomath.ZeroBigDec()
		for _, weight := range normalizedWeights {
			if weight.Asset != tokenIn.Denom {
				swappedTokenIn := sdk.NewCoin(tokenIn.Denom, osmomath.BigDecFromSDKInt(tokenIn.Amount).Mul(weight.Weight).Dec().TruncateInt())
				_, slippage, tokenAmountOut, err = p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, sdk.Coins{swappedTokenIn}, weight.Asset, osmomath.ZeroBigDec())
				if err == nil {
					totalSlippage = totalSlippage.Add(slippage)
					swappedTokenOut := sdk.NewCoin(weight.Asset, tokenAmountOut.Dec().TruncateInt())
					swapsInfos = append(swapsInfos, NewSwapInfo(swappedTokenIn, swappedTokenOut))
				}
			}
		}

		numShares, tokensJoined, err = p.CalcSingleAssetJoinPoolShares(tokensIn, takerFees)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
		}
		poolAssetsByDenom, err := GetPoolAssetsByDenom(p.GetAllPoolAssets())
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
		}
		totalWeight := p.GetBigDecTotalWeight()
		normalizedWeight := poolAssetsByDenom[tokenIn.Denom].GetBigDecWeight().Quo(totalWeight)
		// We multiply the swap fee and taker fees by the normalized weight because it is calculated like this later in CalcSingleAssetJoinPoolShares function
		swapFee = osmomath.OneBigDec().Sub(feeRatio(normalizedWeight, p.PoolParams.GetBigDecSwapFee()))
		takerFee := osmomath.OneBigDec().Sub(feeRatio(normalizedWeight, takerFees))

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, osmomath.BigDec{}, osmomath.BigDec{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
		}
		return tokensJoined, numShares, totalSlippage, osmomath.ZeroBigDec(), swapFee, takerFee, swapsInfos, nil
	}

	initialWeightIn := GetDenomOracleAssetWeight(ctx, oracleKeeper, snapshot.PoolAssets, tokensIn[0].Denom)
	initialWeightOut := osmomath.OneBigDec().Sub(initialWeightIn)

	joinValueWithSlippage, slippage, err := p.CalcJoinValueWithSlippage(ctx, snapshot, oracleKeeper, tokensIn[0], initialWeightOut, params)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
	}

	secondTokenDenom, err := p.GetSecondAssetDenomFromTwoAssetPool(tokensIn[0].Denom)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
	}
	swappedTokenIn := sdk.NewCoin(tokensIn[0].Denom, osmomath.BigDecFromSDKInt(tokensIn[0].Amount).Mul(initialWeightOut).Dec().TruncateInt())
	secondTokenOut, _, _, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, sdk.Coins{swappedTokenIn}, secondTokenDenom, osmomath.ZeroBigDec())
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
	}
	swapsInfos = append(swapsInfos, NewSwapInfo(swappedTokenIn, secondTokenOut))

	tvl, err := p.TVL(ctx, oracleKeeper, accountedPoolKeeper)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
	}

	// Ensure tvl is not zero to avoid division by zero
	if tvl.IsZero() {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, ErrAmountTooLow
	}

	newAssetPools, err := p.NewPoolAssetsAfterSwap(tokensIn, sdk.NewCoins(), snapshot.PoolAssets)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
	}

	weightBalanceBonus, weightBreakingFee, isSwapFee := p.CalculateWeightFees(ctx, oracleKeeper, snapshot.PoolAssets, newAssetPools, tokensIn[0].Denom, params, osmomath.OneBigDec())
	// apply percentage to fees, consider improvement or reduction of other token
	// Other denom weight ratio to reduce the weight breaking fees
	weightBreakingFee = weightBreakingFee.Mul(initialWeightOut)
	weightBalanceBonus = weightBalanceBonus.Mul(initialWeightOut)

	swapFee = osmomath.ZeroBigDec()
	if isSwapFee {
		swapFee = p.GetPoolParams().GetBigDecSwapFee().Mul(initialWeightOut)
	}

	takerFeesFinal = takerFees.Mul(initialWeightOut)

	totalShares := p.GetTotalShares()
	numSharesDec := osmomath.BigDecFromSDKInt(totalShares.Amount).
		Mul(joinValueWithSlippage).Quo(tvl).
		Mul(osmomath.OneBigDec().Sub(weightBreakingFee)).
		Mul(osmomath.OneBigDec().Sub(swapFee.Add(takerFeesFinal)))
	numShares = numSharesDec.Dec().RoundInt()
	err = p.IncreaseLiquidity(numShares, tokensIn)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil, err
	}

	return tokensIn, numShares, slippage, weightBalanceBonus, swapFee, takerFeesFinal, swapsInfos, nil
}

func (p *Pool) GetSecondAssetDenomFromTwoAssetPool(firstAssetDenom string) (string, error) {
	if len(p.PoolAssets) != 2 {
		return "", fmt.Errorf("pool does not have exactly two assets")
	}
	for _, asset := range p.PoolAssets {
		if asset.Token.Denom != firstAssetDenom {
			return asset.Token.Denom, nil
		}
	}
	return "", fmt.Errorf("pool's both asset denoms are same with denom %s", firstAssetDenom)
}
