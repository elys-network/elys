package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcExitValueWithoutSlippage(ctx sdk.Context, oracleKeeper OracleKeeper, accPoolKeeper AccountedPoolKeeper, pool Pool, exitingShares sdkmath.Int, tokenOutDenom string) (sdkmath.LegacyDec, error) {
	tvl, err := pool.TVL(ctx, oracleKeeper, accPoolKeeper)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	totalShares := pool.GetTotalShares()
	refundedShares := sdkmath.LegacyNewDecFromInt(exitingShares)

	// Ensure totalShares is not zero to avoid division by zero
	if totalShares.IsZero() {
		return sdkmath.LegacyZeroDec(), ErrAmountTooLow
	}

	exitValue := tvl.Mul(refundedShares).Quo(sdkmath.LegacyNewDecFromInt(totalShares.Amount))

	if exitingShares.GTE(totalShares.Amount) {
		return sdkmath.LegacyZeroDec(), errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
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
	// 	if exitAmt.LTE(sdkmath.ZeroInt()) {
	// 		continue
	// 	}
	// 	if exitAmt.GTE(asset.Amount) {
	// 		return sdkmath.LegacyZeroDec(), errors.New("too many shares out")
	// 	}
	// 	exitedCoins = exitedCoins.Add(sdk.NewCoin(asset.Denom, exitAmt))
	// }

	// slippageValue := sdkmath.LegacyZeroDec()
	// for _, exitedCoin := range exitedCoins {
	// 	if exitedCoin.Denom == tokenOutDenom {
	// 		continue
	// 	}
	// 	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, exitedCoin.Denom)
	// 	if inTokenPrice.IsZero() {
	// 		return sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", exitedCoin.Denom)
	// 	}
	// 	resizedAmount := sdkmath.LegacyNewDecFromInt(exitedCoin.Amount).
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
	// 		return sdkmath.LegacyZeroDec(), err
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
	exitingShares sdkmath.Int,
	tokenOutDenom string,
	params Params,
	applyWeightBreakingFee bool,
) (exitCoins sdk.Coins, weightBalanceBonus sdkmath.LegacyDec, swapFee sdkmath.LegacyDec, err error) {
	totalShares := pool.GetTotalShares()
	if exitingShares.GTE(totalShares.Amount) {
		return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
	}

	// refundedShares = exitingShares * (1 - exit fee)
	// with 0 exit fee optimization
	refundedShares := exitingShares.ToLegacyDec()

	shareOutRatio := refundedShares.QuoInt(totalShares.Amount)
	// exitedCoins = shareOutRatio * pool liquidity
	exitedCoins := sdk.Coins{}
	poolLiquidity := pool.GetTotalPoolLiquidity()

	if pool.PoolParams.UseOracle && tokenOutDenom != "" {

		accountedAssets := pool.GetAccountedBalance(ctx, accountedPoolKeeper, pool.PoolAssets)
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
		exitValueWithoutSlippage, err := CalcExitValueWithoutSlippage(ctx, oracleKeeper, accountedPoolKeeper, pool, exitingShares, tokenOutDenom)
		if err != nil {
			return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		// Ensure tokenPrice is not zero to avoid division by zero
		if tokenPrice.IsZero() {
			return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrAmountTooLow
		}

		oracleOutAmount := exitValueWithoutSlippage.Quo(tokenPrice)

		tokenOutAmount := oracleOutAmount.RoundInt()
		weightBalanceBonus = sdkmath.LegacyZeroDec()
		isSwapFee := true

		if applyWeightBreakingFee {
			newAssetPools, err := pool.NewPoolAssetsAfterSwap(ctx,
				sdk.Coins{},
				sdk.Coins{sdk.NewCoin(tokenOutDenom, oracleOutAmount.RoundInt())}, accountedAssets,
			)
			if err != nil {
				return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
			}
			var tokenInDenom string
			for _, asset := range newAssetPools {
				if asset.Token.Amount.IsNegative() {
					return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), errors.New("out amount exceeds liquidity balance")
				}

				// As we have two asset pool so other asset will be tokenIn
				if asset.Token.Denom != tokenOutDenom {
					tokenInDenom = asset.Token.Denom
				}
			}

			var weightBreakingFee sdkmath.LegacyDec
			weightBalanceBonus, weightBreakingFee, isSwapFee = pool.CalculateWeightFees(ctx, oracleKeeper, accountedAssets, newAssetPools, tokenInDenom, params, sdkmath.LegacyOneDec())
			// apply percentage to fees, consider improvement or reduction of other token
			// Other denom weight ratio to reduce the weight breaking fees
			initialWeightOut := GetDenomOracleAssetWeight(ctx, pool.PoolId, oracleKeeper, accountedAssets, tokenOutDenom)
			initialWeightIn := sdkmath.LegacyOneDec().Sub(initialWeightOut)
			weightBreakingFee = weightBreakingFee.Mul(initialWeightIn)
			weightBalanceBonus = weightBalanceBonus.Mul(initialWeightIn)

			swapFee := sdkmath.LegacyZeroDec()
			if isSwapFee {
				swapFee = pool.GetPoolParams().SwapFee.Mul(initialWeightIn)
			}

			tokenOutAmount = oracleOutAmount.
				Mul(sdkmath.LegacyOneDec().Sub(weightBreakingFee)).
				Mul(sdkmath.LegacyOneDec().Sub(swapFee)).RoundInt()
		}

		return sdk.Coins{sdk.NewCoin(tokenOutDenom, tokenOutAmount)}, weightBalanceBonus, swapFee, nil
	}

	for _, asset := range poolLiquidity {
		// round down here, due to not wanting to over-exit
		exitAmt := shareOutRatio.MulInt(asset.Amount).TruncateInt()
		if exitAmt.LTE(sdkmath.ZeroInt()) {
			continue
		}
		if exitAmt.GTE(asset.Amount) {
			return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), errors.New("too many shares out")
		}
		exitedCoins = exitedCoins.Add(sdk.NewCoin(asset.Denom, exitAmt))
	}

	return exitedCoins, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), nil
}
