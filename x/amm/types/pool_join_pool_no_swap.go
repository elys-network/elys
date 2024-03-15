package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PoolAssetUSDValue struct {
	Asset string
	Value sdk.Dec
}

type InternalSwapRequest struct {
	InAmount sdk.Coin
	OutToken string
}

func (p *Pool) CalcJoinValueWithoutSlippage(ctx sdk.Context, oracleKeeper OracleKeeper, accountedPoolKeeper AccountedPoolKeeper, tokensIn sdk.Coins) (math.LegacyDec, error) {
	joinValue := sdk.ZeroDec()
	for _, asset := range tokensIn {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Denom)
		if tokenPrice.IsZero() {
			return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", asset.Denom)
		}
		v := tokenPrice.Mul(sdk.NewDecFromInt(asset.Amount))
		joinValue = joinValue.Add(v)
	}
	return joinValue, nil

	// Note: Disable slippage handling for oracle pool due to 1 hour lockup on oracle lp
	// // weights := NormalizedWeights(p.PoolAssets)
	// weights, err := OraclePoolNormalizedWeights(ctx, oracleKeeper, p.PoolAssets)
	// if err != nil {
	// 	return sdk.ZeroDec(), err
	// }

	// inAmounts := []PoolAssetUSDValue{}
	// outAmounts := []PoolAssetUSDValue{}

	// for _, weight := range weights {
	// 	targetAmount := joinValue.Mul(weight.Weight)
	// 	tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, weight.Asset)
	// 	if tokenPrice.IsZero() {
	// 		return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", weight.Asset)
	// 	}
	// 	inAmount := tokenPrice.Mul(sdk.NewDecFromInt(tokensIn.AmountOf(weight.Asset)))
	// 	if targetAmount.GT(inAmount) {
	// 		outAmounts = append(outAmounts, PoolAssetUSDValue{
	// 			Asset: weight.Asset,
	// 			Value: targetAmount.Sub(inAmount),
	// 		})
	// 	}

	// 	if targetAmount.LT(inAmount) {
	// 		inAmounts = append(inAmounts, PoolAssetUSDValue{
	// 			Asset: weight.Asset,
	// 			Value: inAmount.Sub(targetAmount),
	// 		})
	// 	}
	// }

	// internalSwapRequests := []InternalSwapRequest{}
	// for i, j := 0, 0; i < len(inAmounts) && j < len(outAmounts); {
	// 	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, inAmounts[i].Asset)
	// 	if inTokenPrice.IsZero() {
	// 		return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", inAmounts[i].Asset)
	// 	}
	// 	inAsset := inAmounts[i].Asset
	// 	outAsset := outAmounts[j].Asset
	// 	inAmount := sdk.ZeroInt()
	// 	if inAmounts[i].Value.GT(outAmounts[j].Value) {
	// 		inAmount = outAmounts[j].Value.Quo(inTokenPrice).RoundInt()
	// 		j++
	// 	} else if inAmounts[i].Value.LT(outAmounts[j].Value) {
	// 		inAmount = inAmounts[i].Value.Quo(inTokenPrice).RoundInt()
	// 		i++
	// 	} else {
	// 		inAmount = inAmounts[i].Value.Quo(inTokenPrice).RoundInt()
	// 		i++
	// 		j++
	// 	}
	// 	internalSwapRequests = append(internalSwapRequests, InternalSwapRequest{
	// 		InAmount: sdk.NewCoin(inAsset, inAmount),
	// 		OutToken: outAsset,
	// 	})
	// }

	// slippageValue := sdk.ZeroDec()
	// for _, req := range internalSwapRequests {
	// 	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, req.InAmount.Denom)
	// 	if inTokenPrice.IsZero() {
	// 		return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", req.InAmount.Denom)
	// 	}
	// 	resizedAmount := sdk.NewDecFromInt(req.InAmount.Amount).
	// 		Quo(p.PoolParams.ExternalLiquidityRatio).RoundInt()
	// 	slippageAmount, err := p.CalcGivenInSlippage(
	// 		ctx,
	// 		oracleKeeper,
	// 		p,
	// 		sdk.Coins{sdk.NewCoin(req.InAmount.Denom, resizedAmount)},
	// 		req.OutToken,
	// 		accountedPoolKeeper,
	// 	)
	// 	if err != nil {
	// 		return sdk.ZeroDec(), err
	// 	}

	// 	slippageValue = slippageValue.Add(slippageAmount.Mul(inTokenPrice))
	// }
	// joinValueWithoutSlippage := joinValue.Sub(slippageValue)

	// return joinValueWithoutSlippage, nil
}

// JoinPool calculates the number of shares needed for an all-asset join given tokensIn with swapFee applied.
// It updates the liquidity if the pool is joined successfully. If not, returns error.
func (p *Pool) JoinPool(
	ctx sdk.Context, snapshot *Pool,
	oracleKeeper OracleKeeper,
	accountedPoolKeeper AccountedPoolKeeper, tokensIn sdk.Coins,
) (numShares math.Int, slippage sdk.Dec, weightBalanceBonus sdk.Dec, err error) {
	// if it's not single sided liquidity, add at pool ratio
	if len(tokensIn) != 1 {
		numShares, tokensJoined, err := p.CalcJoinPoolNoSwapShares(tokensIn)
		if err != nil {
			return math.Int{}, sdk.ZeroDec(), sdk.ZeroDec(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		p.IncreaseLiquidity(numShares, tokensJoined)
		return numShares, sdk.ZeroDec(), sdk.ZeroDec(), nil
	}

	if !p.PoolParams.UseOracle {
		tokenIn := tokensIn[0]
		totalSlippage := sdk.ZeroDec()
		normalizedWeights := NormalizedWeights(p.PoolAssets)
		for _, weight := range normalizedWeights {
			if weight.Asset != tokenIn.Denom {
				_, slippage, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, weight.Asset, sdk.ZeroDec(), accountedPoolKeeper)
				if err == nil {
					totalSlippage = totalSlippage.Add(slippage.Mul(weight.Weight))
				}
			}
		}

		numShares, tokensJoined, err := p.CalcSingleAssetJoinPoolShares(tokensIn)
		if err != nil {
			return math.Int{}, totalSlippage, sdk.ZeroDec(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		p.IncreaseLiquidity(numShares, tokensJoined)
		return numShares, sdk.ZeroDec(), sdk.ZeroDec(), nil
	}

	joinValueWithoutSlippage, err := p.CalcJoinValueWithoutSlippage(ctx, oracleKeeper, accountedPoolKeeper, tokensIn)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroDec(), sdk.ZeroDec(), err
	}

	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, p.PoolAssets)
	tvl, err := p.TVL(ctx, oracleKeeper)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroDec(), sdk.ZeroDec(), err
	}

	// Ensure tvl is not zero to avoid division by zero
	if tvl.IsZero() {
		return sdk.ZeroInt(), sdk.ZeroDec(), sdk.ZeroDec(), ErrAmountTooLow
	}

	newAssetPools, err := p.NewPoolAssetsAfterSwap(tokensIn, sdk.Coins{})
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroDec(), sdk.ZeroDec(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)

	distanceDiff := weightDistance.Sub(initialWeightDistance)
	weightBreakingFee := sdk.ZeroDec()
	if distanceDiff.IsPositive() {
		// old weight breaking fee implementation
		// weightBreakingFee = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff)

		// we only allow
		tokenInDenom := tokensIn[0].Denom
		// target weight
		targetWeightIn := NormalizedWeight(ctx, p.PoolAssets, tokenInDenom)
		targetWeightOut := sdk.OneDec().Sub(targetWeightIn)

		// weight breaking fee as in Plasma pool
		weightIn := OracleAssetWeight(ctx, oracleKeeper, newAssetPools, tokenInDenom)
		weightOut := sdk.OneDec().Sub(weightIn)
		weightBreakingFee = GetWeightBreakingFee(weightIn, weightOut, targetWeightIn, targetWeightOut, p.PoolParams)

	}
	weightBalanceBonus = weightBreakingFee.Neg()
	if initialWeightDistance.GT(p.PoolParams.ThresholdWeightDifference) && distanceDiff.IsNegative() {
		weightBalanceBonus = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff).Abs()
	}

	totalShares := p.GetTotalShares()
	numSharesDec := sdk.NewDecFromInt(totalShares.Amount).
		Mul(joinValueWithoutSlippage).Quo(tvl).
		Mul(sdk.OneDec().Add(weightBalanceBonus))
	numShares = numSharesDec.RoundInt()
	p.IncreaseLiquidity(numShares, tokensIn)

	// No slippage in oracle pool due to 1 hr lock
	return numShares, sdk.ZeroDec(), weightBalanceBonus, nil
}
