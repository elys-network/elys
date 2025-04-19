package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalcOutAmtGivenIn calculates tokens to be swapped out given the provided
// amount and fee deducted, using solveConstantFunctionInvariant.
func (p Pool) CalcOutAmtGivenIn(
	ctx sdk.Context,
	oracle OracleKeeper,
	snapshot *Pool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee sdkmath.LegacyDec,
	accountedPool AccountedPoolKeeper,
) (sdk.Coin, sdkmath.LegacyDec, error) {
	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
	}

	fmt.Println("-----------CalcOutAmtGivenIn--------------")

	fmt.Println("tokenIn", tokenIn.String())

	tokenAmountInAfterFee := sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Mul(sdkmath.LegacyOneDec().Sub(swapFee))
	fmt.Println("tokenAmountInAfterFee", tokenAmountInAfterFee.String())
	poolTokenInBalance := sdkmath.LegacyNewDecFromInt(poolAssetIn.Token.Amount)
	// accounted pool balance
	acountedPoolAssetInAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetIn.Token.Denom)
	if acountedPoolAssetInAmt.IsPositive() {
		poolTokenInBalance = sdkmath.LegacyNewDecFromInt(acountedPoolAssetInAmt)
	}

	poolTokenOutBalance := sdkmath.LegacyNewDecFromInt(poolAssetOut.Token.Amount)
	// accounted pool balance
	accountedPoolAssetOutAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetOut.Token.Denom)
	if accountedPoolAssetOutAmt.IsPositive() {
		poolTokenOutBalance = sdkmath.LegacyNewDecFromInt(accountedPoolAssetOutAmt)
	}

	fmt.Println("poolTokenInBalance", poolTokenInBalance.String())
	fmt.Println("poolTokenOutBalance", poolTokenOutBalance.String())

	poolPostSwapInBalance := poolTokenInBalance.Add(tokenAmountInAfterFee)

	outWeight := sdkmath.LegacyNewDecFromInt(poolAssetOut.Weight)
	inWeight := sdkmath.LegacyNewDecFromInt(poolAssetIn.Weight)
	if p.PoolParams.UseOracle {
		_, poolAssetIn, poolAssetOut, err := snapshot.parsePoolAssets(tokensIn, tokenOutDenom)
		if err != nil {
			return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
		}
		fmt.Println("-----_Weight Calculation_-----")

		fmt.Println("poolAssetIn", poolAssetIn.String())
		fmt.Println("poolAssetOut", poolAssetOut.String())
		oraclePoolWeights := []AssetWeight{}

		totalWeight := sdkmath.LegacyZeroDec()

		for _, asset := range []PoolAsset{poolAssetIn, poolAssetOut} {
			fmt.Println("denom: ", asset.Token.Denom, " amount: ", asset.Token.Amount.String())
			tokenPrice := oracle.GetAssetPriceFromDenom(ctx, asset.Token.Denom)
			amount := asset.Token.Amount
			weight := amount.ToLegacyDec().Mul(tokenPrice)
			oraclePoolWeights = append(oraclePoolWeights, AssetWeight{
				Asset:  asset.Token.Denom,
				Weight: weight,
			})
			totalWeight = totalWeight.Add(weight)
		}
		fmt.Println("totalWeight", totalWeight.String())
		for i, asset := range oraclePoolWeights {
			oraclePoolWeights[i].Weight = asset.Weight.Quo(totalWeight)
			fmt.Println("denom: ", asset.Asset, " weight: ", asset.Weight.String())
		}

		oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, p.PoolId, oracle, []PoolAsset{poolAssetIn, poolAssetOut})
		if err != nil {
			return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
		}
		inWeight = oracleWeights[0].Weight
		outWeight = oracleWeights[1].Weight
		fmt.Println("-----_Weight Calculation_-----")
	}

	fmt.Println("inWeight", inWeight.String())
	fmt.Println("outWeight", outWeight.String())

	// deduct swapfee on the tokensIn
	// delta balanceOut is positive(tokens inside the pool decreases)
	tokenAmountOut, err := solveConstantFunctionInvariant(
		poolTokenInBalance,
		poolPostSwapInBalance,
		inWeight,
		poolTokenOutBalance,
		outWeight,
	)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
	}

	if tokenAmountOut.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), ErrTokenOutAmountZero
	}

	fmt.Println("tokenAmountOut", tokenAmountOut.String())
	rate, err := p.GetTokenARate(ctx, oracle, snapshot, tokenIn.Denom, tokenOutDenom, accountedPool)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
	}

	amountOutWithoutSlippage := tokenAmountInAfterFee.Mul(rate)

	// check if amountOutWithoutSlippage is zero to avoid division by zero
	if amountOutWithoutSlippage.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), errorsmod.Wrapf(ErrInvalidMathApprox, "amount out without slippage must be positive")
	}

	slippage := sdkmath.LegacyOneDec().Sub(tokenAmountOut.Quo(amountOutWithoutSlippage))

	fmt.Println("slippage", slippage.String())
	// We ignore the decimal component, as we round down the token amount out.
	tokenAmountOutInt := tokenAmountOut.TruncateInt()
	if !tokenAmountOutInt.IsPositive() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), ErrTokenOutAmountZero
	}

	fmt.Println("tokenAmountOutInt", tokenAmountOutInt.String())
	fmt.Println("-----------CalcOutAmtGivenIn--------------")
	return sdk.NewCoin(tokenOutDenom, tokenAmountOutInt), slippage, nil
}
