package types

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func CalcExitValueWithSlippage(ctx sdk.Context, oracleKeeper OracleKeeper, accPoolKeeper AccountedPoolKeeper,
	pool Pool, exitingShares sdkmath.Int, tokenOutDenom string,
	weightMultiplier osmomath.BigDec, applyFee bool, params Params) (osmomath.BigDec, osmomath.BigDec, sdk.Coins, error) {
	tvl, err := pool.TVL(ctx, oracleKeeper, accPoolKeeper)
	if err != nil {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, err
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
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, fmt.Errorf("token in denom not found")
	}

	totalShares := pool.GetTotalShares()
	refundedShares := osmomath.BigDecFromSDKInt(exitingShares)

	// Ensure totalShares is not zero to avoid division by zero
	if totalShares.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, ErrAmountTooLow
	}

	exitValue := tvl.Mul(refundedShares).Quo(osmomath.BigDecFromSDKInt(totalShares.Amount))

	if !applyFee {
		return exitValue, osmomath.ZeroBigDec(), sdk.Coins{}, nil
	}

	inTokenPrice := oracleKeeper.GetDenomPrice(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, fmt.Errorf("token price not set: %s", tokenInDenom)
	}

	outTokenPrice := oracleKeeper.GetDenomPrice(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, fmt.Errorf("token price not set: %s", tokenOutDenom)
	}

	externalLiquidityRatio, err := pool.GetAssetExternalLiquidityRatio(tokenOutDenom)
	if err != nil {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, err
	}
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.LT(osmomath.OneBigDec()) {
		externalLiquidityRatio = osmomath.OneBigDec()
	}

	// tokenIn amount will be
	tokenInAmount := exitValue.Quo(inTokenPrice)
	weightedAmount := tokenInAmount.Mul(weightMultiplier)
	resizedAmount := osmomath.BigDecFromSDKInt(weightedAmount.Dec().TruncateInt()).
		Quo(externalLiquidityRatio).Dec().RoundInt()
	slippageAmount, err := pool.CalcGivenInSlippage(
		ctx,
		oracleKeeper,
		&pool,
		sdk.Coins{sdk.NewCoin(tokenInDenom, resizedAmount)},
		tokenOutDenom,
		accPoolKeeper,
	)
	if err != nil {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, err
	}
	slippageAmount = slippageAmount.Mul(externalLiquidityRatio)

	slippageValue := slippageAmount.Mul(outTokenPrice)
	slippage := slippageValue.Quo(exitValue)

	minSlippage := params.GetBigDecMinSlippage().Mul(weightMultiplier)
	if slippage.LT(minSlippage) {
		slippage = minSlippage
		slippageValue = exitValue.Mul(minSlippage)
	}

	exitValueWithSlippage := exitValue.Sub(slippageValue)

	if exitingShares.GTE(totalShares.Amount) {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
	}

	slippageCoins := sdk.Coins{sdk.NewCoin(tokenOutDenom, slippageAmount.Dec().TruncateInt())}

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
	takerFees osmomath.BigDec,
	applyFee bool,
) (exitCoins sdk.Coins, weightBalanceBonus osmomath.BigDec, slippage osmomath.BigDec, swapFee osmomath.BigDec, takerFeesFinal osmomath.BigDec, slippageCoins sdk.Coins, err error) {
	totalShares := pool.GetTotalShares()
	if exitingShares.GTE(totalShares.Amount) {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
	}

	// refundedShares = exitingShares * (1 - exit fee)
	// with 0 exit fee optimization
	refundedShares := osmomath.BigDecFromSDKInt(exitingShares)

	shareOutRatio := refundedShares.Quo(osmomath.BigDecFromSDKInt(totalShares.Amount))
	// exitedCoins = shareOutRatio * pool liquidity
	exitedCoins := sdk.Coins{}
	poolLiquidity := pool.GetTotalPoolLiquidity()

	if pool.PoolParams.UseOracle && tokenOutDenom != "" {

		accountedAssets := pool.GetAccountedBalance(ctx, accountedPoolKeeper, pool.PoolAssets)
		tokenPrice := oracleKeeper.GetDenomPrice(ctx, tokenOutDenom)

		initialWeightOut := GetDenomOracleAssetWeight(ctx, pool.PoolId, oracleKeeper, accountedAssets, tokenOutDenom)
		initialWeightIn := osmomath.OneBigDec().Sub(initialWeightOut)

		exitValueWithSlippage, slippage, slippageCoins, err := CalcExitValueWithSlippage(ctx, oracleKeeper, accountedPoolKeeper, pool, exitingShares, tokenOutDenom, initialWeightIn, applyFee, params)
		if err != nil {
			return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, err
		}

		// Ensure tokenPrice is not zero to avoid division by zero
		if tokenPrice.IsZero() {
			return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, ErrAmountTooLow
		}

		oracleOutAmount := exitValueWithSlippage.Quo(tokenPrice)

		tokenOutAmount := oracleOutAmount.Dec().RoundInt()
		weightBalanceBonus = osmomath.ZeroBigDec()
		takerFeesFinal = osmomath.ZeroBigDec()
		isSwapFee := true
		swapFee = osmomath.ZeroBigDec()

		if applyFee {
			newAssetPools, err := pool.NewPoolAssetsAfterSwap(ctx,
				sdk.Coins{},
				sdk.Coins{sdk.NewCoin(tokenOutDenom, oracleOutAmount.Dec().RoundInt())}, accountedAssets,
			)
			if err != nil {
				return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, err
			}
			var tokenInDenom string
			for _, asset := range newAssetPools {
				if asset.Token.Amount.IsNegative() {
					return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, errors.New("out amount exceeds liquidity balance")
				}

				// As we have two asset pool so other asset will be tokenIn
				if asset.Token.Denom != tokenOutDenom {
					tokenInDenom = asset.Token.Denom
				}
			}

			var weightBreakingFee osmomath.BigDec
			weightBalanceBonus, weightBreakingFee, isSwapFee = pool.CalculateWeightFees(ctx, oracleKeeper, accountedAssets, newAssetPools, tokenInDenom, params, osmomath.OneBigDec())
			// apply percentage to fees, consider improvement or reduction of other token
			// Other denom weight ratio to reduce the weight breaking fees
			weightBreakingFee = weightBreakingFee.Mul(initialWeightIn)
			weightBalanceBonus = weightBalanceBonus.Mul(initialWeightIn)

			if isSwapFee {
				swapFee = pool.GetPoolParams().GetBigDecSwapFee().Mul(initialWeightIn)
			}

			takerFeesFinal = takerFees.Mul(initialWeightIn)

			tokenOutAmount = (oracleOutAmount.
				Mul(osmomath.OneBigDec().Sub(weightBreakingFee)).
				Mul(osmomath.OneBigDec().Sub(swapFee.Add(takerFeesFinal)))).Dec().RoundInt()
		}

		return sdk.Coins{sdk.NewCoin(tokenOutDenom, tokenOutAmount)}, weightBalanceBonus, slippage, swapFee, takerFeesFinal, slippageCoins, nil
	}

	for _, asset := range poolLiquidity {
		// round down here, due to not wanting to over-exit
		exitAmt := shareOutRatio.Mul(osmomath.BigDecFromSDKInt(asset.Amount)).Dec().TruncateInt()
		if exitAmt.LTE(sdkmath.ZeroInt()) {
			continue
		}
		if exitAmt.GTE(asset.Amount) {
			return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, errors.New("too many shares out")
		}
		exitedCoins = exitedCoins.Add(sdk.NewCoin(asset.Denom, exitAmt))
	}

	return exitedCoins, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, nil
}
