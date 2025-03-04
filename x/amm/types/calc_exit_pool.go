package types

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcExitValueWithSlippage(ctx sdk.Context, oracleKeeper OracleKeeper, accPoolKeeper AccountedPoolKeeper,
	pool Pool, exitingShares sdkmath.Int, tokenOutDenom string,
	weightMultiplier sdkmath.LegacyDec, applyFee bool, params Params) (sdkmath.LegacyDec, sdkmath.LegacyDec, sdk.Coins, error) {
	tvl, err := pool.TVL(ctx, oracleKeeper, accPoolKeeper)
	if err != nil {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, err
	}

	// As this is 2 token pool, tokenOut will be
	tokenInDenom := ""
	for _, asset := range pool.PoolAssets {
		if asset.Token.Denom == tokenOutDenom {
			continue
		}
		tokenInDenom = asset.Token.Denom
	}
	// Not possible, but we might require this when we have pools with assets more than 2
	if tokenInDenom == "" {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, fmt.Errorf("token in denom not found")
	}

	totalShares := pool.GetTotalShares()
	refundedShares := sdkmath.LegacyNewDecFromInt(exitingShares)

	// Ensure totalShares is not zero to avoid division by zero
	if totalShares.IsZero() {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, ErrAmountTooLow
	}

	exitValue := tvl.Mul(refundedShares).Quo(sdkmath.LegacyNewDecFromInt(totalShares.Amount))

	if !applyFee {
		return exitValue, sdkmath.LegacyZeroDec(), sdk.Coins{}, nil
	}

	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, fmt.Errorf("token price not set: %s", tokenInDenom)
	}

	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, fmt.Errorf("token price not set: %s", tokenOutDenom)
	}

	externalLiquidityRatio, err := pool.GetAssetExternalLiquidityRatio(tokenOutDenom)
	if err != nil {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, err
	}
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.LT(sdkmath.LegacyOneDec()) {
		externalLiquidityRatio = sdkmath.LegacyOneDec()
	}

	// tokenIn amount will be
	tokenInAmount := exitValue.Quo(inTokenPrice)
	weightedAmount := tokenInAmount.Mul(weightMultiplier)
	resizedAmount := sdkmath.LegacyNewDecFromInt(weightedAmount.TruncateInt()).
		Quo(externalLiquidityRatio).RoundInt()
	slippageAmount, err := pool.CalcGivenInSlippage(
		ctx,
		oracleKeeper,
		&pool,
		sdk.Coins{sdk.NewCoin(tokenInDenom, resizedAmount)},
		tokenOutDenom,
		accPoolKeeper,
	)
	if err != nil {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, err
	}
	slippageAmount = slippageAmount.Mul(externalLiquidityRatio)

	slippageValue := slippageAmount.Mul(outTokenPrice)
	slippage := slippageValue.Quo(exitValue)

	minSlippage := params.MinSlippage.Mul(weightMultiplier)
	if slippage.LT(minSlippage) {
		slippage = minSlippage
		slippageValue = exitValue.Mul(minSlippage)
	}

	exitValueWithSlippage := exitValue.Sub(slippageValue)

	if exitingShares.GTE(totalShares.Amount) {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
	}

	slippageCoins := sdk.Coins{sdk.NewCoin(tokenOutDenom, slippageAmount.TruncateInt())}

	return exitValueWithSlippage, slippage, slippageCoins, nil
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
	takerFees sdkmath.LegacyDec,
	applyFee bool,
) (exitCoins sdk.Coins, weightBalanceBonus sdkmath.LegacyDec, slippage sdkmath.LegacyDec, swapFee sdkmath.LegacyDec, takerFeesFinal sdkmath.LegacyDec, slippageCoins sdk.Coins, err error) {
	totalShares := pool.GetTotalShares()
	if exitingShares.GTE(totalShares.Amount) {
		return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
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

		initialWeightOut := GetDenomOracleAssetWeight(ctx, pool.PoolId, oracleKeeper, accountedAssets, tokenOutDenom)
		initialWeightIn := sdkmath.LegacyOneDec().Sub(initialWeightOut)

		exitValueWithSlippage, slippage, slippageCoins, err := CalcExitValueWithSlippage(ctx, oracleKeeper, accountedPoolKeeper, pool, exitingShares, tokenOutDenom, initialWeightIn, applyFee, params)
		if err != nil {
			return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, err
		}

		// Ensure tokenPrice is not zero to avoid division by zero
		if tokenPrice.IsZero() {
			return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, ErrAmountTooLow
		}

		oracleOutAmount := exitValueWithSlippage.Quo(tokenPrice)

		tokenOutAmount := oracleOutAmount.RoundInt()
		weightBalanceBonus = sdkmath.LegacyZeroDec()
		takerFeesFinal = sdkmath.LegacyZeroDec()
		isSwapFee := true
		swapFee = sdkmath.LegacyZeroDec()

		if applyFee {
			newAssetPools, err := pool.NewPoolAssetsAfterSwap(ctx,
				sdk.Coins{},
				sdk.Coins{sdk.NewCoin(tokenOutDenom, oracleOutAmount.RoundInt())}, accountedAssets,
			)
			if err != nil {
				return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, err
			}
			var tokenInDenom string
			for _, asset := range newAssetPools {
				if asset.Token.Amount.IsNegative() {
					return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, errors.New("out amount exceeds liquidity balance")
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
			weightBreakingFee = weightBreakingFee.Mul(initialWeightIn)
			weightBalanceBonus = weightBalanceBonus.Mul(initialWeightIn)

			if isSwapFee {
				swapFee = pool.GetPoolParams().SwapFee.Mul(initialWeightIn)
			}

			takerFeesFinal = takerFees.Mul(initialWeightIn)

			tokenOutAmount = (oracleOutAmount.
				Mul(sdkmath.LegacyOneDec().Sub(weightBreakingFee)).
				Mul(sdkmath.LegacyOneDec().Sub(swapFee.Add(takerFeesFinal)))).RoundInt()
		}

		return sdk.Coins{sdk.NewCoin(tokenOutDenom, tokenOutAmount)}, weightBalanceBonus, slippage, swapFee, takerFeesFinal, slippageCoins, nil
	}

	for _, asset := range poolLiquidity {
		// round down here, due to not wanting to over-exit
		exitAmt := shareOutRatio.MulInt(asset.Amount).TruncateInt()
		if exitAmt.LTE(sdkmath.ZeroInt()) {
			continue
		}
		if exitAmt.GTE(asset.Amount) {
			return sdk.Coins{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, errors.New("too many shares out")
		}
		exitedCoins = exitedCoins.Add(sdk.NewCoin(asset.Denom, exitAmt))
	}

	return exitedCoins, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdk.Coins{}, nil
}
