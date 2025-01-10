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

func (p *Pool) CalcJoinValueWithoutSlippage(ctx sdk.Context, oracleKeeper OracleKeeper, accountedPoolKeeper AccountedPoolKeeper, tokensIn sdk.Coins) (elystypes.Dec34, error) {
	joinValue := elystypes.ZeroDec34()
	for _, asset := range tokensIn {
		tokenPrice, decimals := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Denom)
		if tokenPrice.IsZero() {
			return elystypes.ZeroDec34(), fmt.Errorf("token price not set: %s", asset.Denom)
		}
		v := tokenPrice.MulInt(asset.Amount).QuoInt(OneTokenUnit(decimals))
		joinValue = joinValue.Add(v)
	}
	return joinValue, nil
}

// JoinPool calculates the number of shares needed for an all-asset join given tokensIn with swapFee applied.
// It updates the liquidity if the pool is joined successfully. If not, returns error.
func (p *Pool) JoinPool(
	ctx sdk.Context, snapshot *Pool,
	oracleKeeper OracleKeeper,
	accountedPoolKeeper AccountedPoolKeeper, tokensIn sdk.Coins,
	params Params,
) (tokensJoined sdk.Coins, numShares sdkmath.Int, slippage elystypes.Dec34, weightBalanceBonus elystypes.Dec34, err error) {
	// if it's not single sided liquidity, add at pool ratio
	if len(tokensIn) != 1 {
		numShares, tokensJoined, err := p.CalcJoinPoolNoSwapShares(tokensIn)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}
		return tokensJoined, numShares, elystypes.ZeroDec34(), elystypes.ZeroDec34(), nil
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
			return sdk.NewCoins(), sdkmath.Int{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}

		// update pool with the calculated share and liquidity needed to join pool
		err = p.IncreaseLiquidity(numShares, tokensJoined)
		if err != nil {
			return sdk.NewCoins(), sdkmath.Int{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
		}
		return tokensJoined, numShares, totalSlippage, elystypes.ZeroDec34(), nil
	}

	joinValueWithoutSlippage, err := p.CalcJoinValueWithoutSlippage(ctx, oracleKeeper, accountedPoolKeeper, tokensIn)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}

	accountedAssets := p.GetAccountedBalance(ctx, accountedPoolKeeper, p.PoolAssets)
	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, accountedAssets)
	tvl, err := p.TVL(ctx, oracleKeeper, accountedPoolKeeper)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}

	// Ensure tvl is not zero to avoid division by zero
	if tvl.IsZero() {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), ErrAmountTooLow
	}

	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx, tokensIn, sdk.NewCoins(), accountedAssets)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// we only allow
	tokenInDenom := tokensIn[0].Denom
	// target weight
	targetWeightIn := GetDenomNormalizedWeight(p.PoolAssets, tokenInDenom)
	targetWeightOut := elystypes.OneDec34().Sub(targetWeightIn)

	// weight breaking fee as in Plasma pool
	finalWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, newAssetPools, tokenInDenom)
	finalWeightOut := elystypes.OneDec34().Sub(finalWeightIn)

	initialAssetPools, err := p.NewPoolAssetsAfterSwap(ctx,
		sdk.NewCoins(),
		sdk.NewCoins(), accountedAssets,
	)
	initialWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, initialAssetPools, tokenInDenom)
	initialWeightOut := elystypes.OneDec34().Sub(initialWeightIn)
	weightBreakingFee := GetWeightBreakingFee(finalWeightIn, finalWeightOut, targetWeightIn, targetWeightOut, initialWeightIn, initialWeightOut, distanceDiff, params)
	// apply percentage to fees, consider improvement or reduction of other token
	// Other denom weight ratio to reduce the weight breaking fees
	weightBreakingFee = weightBreakingFee.Mul(finalWeightOut)

	// weight recovery reward = weight breaking fee * weight breaking fee portion
	weightRecoveryReward := weightBreakingFee.Mul(elystypes.NewDec34FromLegacyDec(params.WeightBreakingFeePortion))

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus = weightBreakingFee.Neg()
	if initialWeightDistance.GT(elystypes.NewDec34FromLegacyDec(params.ThresholdWeightDifference)) && distanceDiff.IsNegative() {
		weightBalanceBonus = weightRecoveryReward
		// set weight breaking fee to zero if bonus is applied
		weightBreakingFee = elystypes.ZeroDec34()
	}

	totalShares := p.GetTotalShares()
	numSharesDec := elystypes.NewDec34FromInt(totalShares.Amount).
		Mul(joinValueWithoutSlippage).Quo(tvl).
		Mul(elystypes.OneDec34().Sub(weightBreakingFee))
	numShares = numSharesDec.ToInt()
	err = p.IncreaseLiquidity(numShares, tokensIn)
	if err != nil {
		return sdk.NewCoins(), sdkmath.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), err
	}

	// No slippage in oracle pool due to 1 hr lock
	return tokensIn, numShares, elystypes.ZeroDec34(), weightBalanceBonus, nil
}
