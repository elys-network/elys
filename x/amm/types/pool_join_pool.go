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

func (p *Pool) CalcJoinValueWithoutSlippage(ctx sdk.Context, oracleKeeper OracleKeeper, accountedPoolKeeper AccountedPoolKeeper, tokensIn sdk.Coins) (sdkmath.LegacyDec, error) {
	joinValue := sdkmath.LegacyZeroDec()
	for _, asset := range tokensIn {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Denom)
		if tokenPrice.IsZero() {
			return sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", asset.Denom)
		}
		v := tokenPrice.Mul(sdkmath.LegacyNewDecFromInt(asset.Amount))
		joinValue = joinValue.Add(v)
	}
	return joinValue, nil

	// Note: Disable slippage handling for oracle pool due to 1 hour lockup on oracle lp
	// // weights := NormalizedWeights(p.PoolAssets)
	// weights, err := GetOraclePoolNormalizedWeights(ctx, oracleKeeper, p.PoolAssets)
	// if err != nil {
	// 	return sdkmath.LegacyZeroDec(), err
	// }

	// inAmounts := []PoolAssetUSDValue{}
	// outAmounts := []PoolAssetUSDValue{}

	// for _, weight := range weights {
	// 	targetAmount := joinValue.Mul(weight.Weight)
	// 	tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, weight.Asset)
	// 	if tokenPrice.IsZero() {
	// 		return sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", weight.Asset)
	// 	}
	// 	inAmount := tokenPrice.Mul(sdkmath.LegacyNewDecFromInt(tokensIn.AmountOf(weight.Asset)))
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
	// 		return sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", inAmounts[i].Asset)
	// 	}
	// 	inAsset := inAmounts[i].Asset
	// 	outAsset := outAmounts[j].Asset
	// 	inAmount := sdkmath.ZeroInt()
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

	// slippageValue := sdkmath.LegacyZeroDec()
	// for _, req := range internalSwapRequests {
	// 	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, req.InAmount.Denom)
	// 	if inTokenPrice.IsZero() {
	// 		return sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", req.InAmount.Denom)
	// 	}
	// 	resizedAmount := sdkmath.LegacyNewDecFromInt(req.InAmount.Amount).
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
	// 		return sdkmath.LegacyZeroDec(), err
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
	params Params,
) (tokensJoined sdk.Coins, numShares sdkmath.Int, slippage sdkmath.LegacyDec, weightBalanceBonus sdkmath.LegacyDec, err error) {
	// if it's not single sided liquidity, add at pool ratio
	if len(tokensIn) != 1 {
		numShares, tokensJoined, err := p.CalcJoinPoolNoSwapShares(tokensIn)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, sdkmath.LegacyDec{}, sdkmath.LegacyDec{}, err
		}
		return tokensJoined, numShares, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), nil
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
			return sdk.NewCoins(), sdkmath.Int{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, sdkmath.LegacyDec{}, sdkmath.LegacyDec{}, err
		}
		return tokensJoined, numShares, totalSlippage, sdkmath.LegacyZeroDec(), nil
	}

	joinValueWithoutSlippage, err := p.CalcJoinValueWithoutSlippage(ctx, oracleKeeper, accountedPoolKeeper, tokensIn)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	accountedAssets := p.GetAccountedBalance(ctx, accountedPoolKeeper, p.PoolAssets)
	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, accountedAssets)
	tvl, err := p.TVL(ctx, oracleKeeper, accountedPoolKeeper)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// Ensure tvl is not zero to avoid division by zero
	if tvl.IsZero() {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrAmountTooLow
	}

	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx, tokensIn, sdk.NewCoins(), accountedAssets)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// we only allow
	tokenInDenom := tokensIn[0].Denom
	// target weight
	targetWeightIn := GetDenomNormalizedWeight(p.PoolAssets, tokenInDenom)
	targetWeightOut := sdkmath.LegacyOneDec().Sub(targetWeightIn)

	// weight breaking fee as in Plasma pool
	finalWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, newAssetPools, tokenInDenom)
	finalWeightOut := sdkmath.LegacyOneDec().Sub(finalWeightIn)

	initialAssetPools, err := p.NewPoolAssetsAfterSwap(ctx,
		sdk.NewCoins(),
		sdk.NewCoins(), accountedAssets,
	)
	initialWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, initialAssetPools, tokenInDenom)
	initialWeightOut := sdkmath.LegacyOneDec().Sub(initialWeightIn)
	weightBreakingFee := GetWeightBreakingFee(finalWeightIn, finalWeightOut, targetWeightIn, targetWeightOut, initialWeightIn, initialWeightOut, distanceDiff, params)
	// apply percentage to fees, consider improvement or reduction of other token
	// Other denom weight ratio to reduce the weight breaking fees
	weightBreakingFee = weightBreakingFee.Mul(finalWeightOut)

	// weight recovery reward = weight breaking fee * weight breaking fee portion
	weightRecoveryReward := weightBreakingFee.Mul(params.WeightBreakingFeePortion)

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus = weightBreakingFee.Neg()
	if initialWeightDistance.GT(params.ThresholdWeightDifference) && distanceDiff.IsNegative() {
		weightBalanceBonus = weightRecoveryReward
		// set weight breaking fee to zero if bonus is applied
		weightBreakingFee = sdkmath.LegacyZeroDec()
	}

	totalShares := p.GetTotalShares()
	numSharesDec := sdkmath.LegacyNewDecFromInt(totalShares.Amount).
		Mul(joinValueWithoutSlippage).Quo(tvl).
		Mul(sdkmath.LegacyOneDec().Sub(weightBreakingFee))
	numShares = numSharesDec.RoundInt()
	err = p.IncreaseLiquidity(numShares, tokensIn)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// No slippage in oracle pool due to 1 hr lock
	return tokensIn, numShares, sdkmath.LegacyZeroDec(), weightBalanceBonus, nil
}
