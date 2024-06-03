package types

import (
	"errors"
	fmt "fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcExitValueWithoutSlippage(ctx sdk.Context, oracleKeeper OracleKeeper, accPoolKeeper AccountedPoolKeeper, pool Pool, exitingShares math.Int, tokenOutDenom string) (sdk.Dec, error) {
	tvl, err := pool.TVL(ctx, oracleKeeper)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	totalShares := pool.GetTotalShares()
	var refundedShares sdk.Dec
	refundedShares = sdk.NewDecFromInt(exitingShares)

	// Ensure totalShares is not zero to avoid division by zero
	if totalShares.IsZero() {
		return sdk.ZeroDec(), ErrAmountTooLow
	}

	exitValue := tvl.Mul(refundedShares).Quo(sdk.NewDecFromInt(totalShares.Amount))

	if exitingShares.GTE(totalShares.Amount) {
		return sdk.ZeroDec(), errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
	}

	return exitValue, nil

	// Note: Disable slippage handling for oracle pool due to 1 hour lockup on oracle lp
	// shareOutRatio := refundedShares.QuoInt(totalShares.Amount)
	// // exitedCoins = shareOutRatio * pool liquidity
	// exitedCoins := sdk.Coins{}
	// poolLiquidity := pool.GetTotalPoolLiquidity()

	// for _, asset := range poolLiquidity {
	// 	// round down here, due to not wanting to over-exit
	// 	exitAmt := shareOutRatio.MulInt(asset.Amount).TruncateInt()
	// 	if exitAmt.LTE(sdk.ZeroInt()) {
	// 		continue
	// 	}
	// 	if exitAmt.GTE(asset.Amount) {
	// 		return sdk.ZeroDec(), errors.New("too many shares out")
	// 	}
	// 	exitedCoins = exitedCoins.Add(sdk.NewCoin(asset.Denom, exitAmt))
	// }

	// slippageValue := sdk.ZeroDec()
	// for _, exitedCoin := range exitedCoins {
	// 	if exitedCoin.Denom == tokenOutDenom {
	// 		continue
	// 	}
	// 	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, exitedCoin.Denom)
	// 	if inTokenPrice.IsZero() {
	// 		return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", exitedCoin.Denom)
	// 	}
	// 	resizedAmount := sdk.NewDecFromInt(exitedCoin.Amount).
	// 		Quo(pool.PoolParams.ExternalLiquidityRatio).RoundInt()
	// 	slippageAmount, err := pool.CalcGivenInSlippage(
	// 		ctx,
	// 		oracleKeeper,
	// 		&pool,
	// 		sdk.Coins{sdk.NewCoin(exitedCoin.Denom, resizedAmount)},
	// 		tokenOutDenom,
	// 		accPoolKeeper,
	// 	)
	// 	if err != nil {
	// 		return sdk.ZeroDec(), err
	// 	}

	// 	slippageValue = slippageValue.Add(slippageAmount.Mul(inTokenPrice))
	// }
	// exitValueWithoutSlippage := exitValue.Sub(slippageValue)
	// return exitValueWithoutSlippage, nil
}

// CalcExitPool returns how many tokens should come out, when exiting k LP shares against a "standard" CFMM
func CalcExitPool(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	pool Pool,
	accountedPoolKeeper AccountedPoolKeeper,
	exitingShares math.Int,
	tokenOutDenom string,
) (exitCoins sdk.Coins, weightBalanceBonus math.LegacyDec, err error) {
	totalShares := pool.GetTotalShares()
	if exitingShares.GTE(totalShares.Amount) {
		return sdk.Coins{}, math.LegacyZeroDec(), errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
	}

	// refundedShares = exitingShares * (1 - exit fee)
	// with 0 exit fee optimization
	var refundedShares sdk.Dec
	refundedShares = sdk.NewDecFromInt(exitingShares)

	shareOutRatio := refundedShares.QuoInt(totalShares.Amount)
	// exitedCoins = shareOutRatio * pool liquidity
	exitedCoins := sdk.Coins{}
	poolLiquidity := pool.GetTotalPoolLiquidity()

	if pool.PoolParams.UseOracle && tokenOutDenom != "" {
		initialWeightDistance := pool.WeightDistanceFromTarget(ctx, oracleKeeper, pool.PoolAssets)
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
		exitValueWithoutSlippage, err := CalcExitValueWithoutSlippage(ctx, oracleKeeper, accountedPoolKeeper, pool, exitingShares, tokenOutDenom)
		if err != nil {
			return sdk.Coins{}, math.LegacyZeroDec(), err
		}

		// Ensure tokenPrice is not zero to avoid division by zero
		if tokenPrice.IsZero() {
			return sdk.Coins{}, math.LegacyZeroDec(), ErrAmountTooLow
		}

		oracleOutAmount := exitValueWithoutSlippage.Quo(tokenPrice)

		newAssetPools, err := pool.NewPoolAssetsAfterSwap(
			sdk.Coins{},
			sdk.Coins{sdk.NewCoin(tokenOutDenom, oracleOutAmount.RoundInt())},
		)
		if err != nil {
			return sdk.Coins{}, math.LegacyZeroDec(), err
		}
		for _, asset := range newAssetPools {
			if asset.Token.Amount.IsNegative() {
				return sdk.Coins{}, math.LegacyZeroDec(), fmt.Errorf("out amount exceeds liquidity balance")
			}
		}

		weightDistance := pool.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
		distanceDiff := weightDistance.Sub(initialWeightDistance)
		weightBreakingFee := sdk.ZeroDec()
		if distanceDiff.IsPositive() {
			// old weight breaking fee implementation
			// weightBreakingFee = pool.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff)

			// target weight
			targetWeightOut := NormalizedWeight(ctx, pool.PoolAssets, tokenOutDenom)
			targetWeightIn := sdk.OneDec().Sub(targetWeightOut)

			// weight breaking fee as in Plasma pool
			weightOut := OracleAssetWeight(ctx, oracleKeeper, newAssetPools, tokenOutDenom)
			weightIn := sdk.OneDec().Sub(weightOut)

			weightBreakingFee = GetWeightBreakingFee(weightIn, weightOut, targetWeightIn, targetWeightOut, pool.PoolParams)
		}

		tokenOutAmount := oracleOutAmount.Mul(sdk.OneDec().Sub(weightBreakingFee)).RoundInt()
		return sdk.Coins{sdk.NewCoin(tokenOutDenom, tokenOutAmount)}, weightBreakingFee.Neg(), nil
	}

	for _, asset := range poolLiquidity {
		// round down here, due to not wanting to over-exit
		exitAmt := shareOutRatio.MulInt(asset.Amount).TruncateInt()
		if exitAmt.LTE(sdk.ZeroInt()) {
			continue
		}
		if exitAmt.GTE(asset.Amount) {
			return sdk.Coins{}, math.LegacyZeroDec(), errors.New("too many shares out")
		}
		exitedCoins = exitedCoins.Add(sdk.NewCoin(asset.Denom, exitAmt))
	}

	return exitedCoins, math.LegacyZeroDec(), nil
}
