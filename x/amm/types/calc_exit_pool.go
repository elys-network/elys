package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
)

func CalcExitValueWithoutSlippage(ctx sdk.Context, oracleKeeper OracleKeeper, accPoolKeeper AccountedPoolKeeper, pool Pool, exitingShares sdkmath.Int, tokenOutDenom string) (elystypes.Dec34, error) {
	tvl, err := pool.TVL(ctx, oracleKeeper, accPoolKeeper)
	if err != nil {
		return elystypes.ZeroDec34(), err
	}

	totalShares := pool.GetTotalShares()
	refundedShares := elystypes.NewDec34FromInt(exitingShares)

	// Ensure totalShares is not zero to avoid division by zero
	if totalShares.IsZero() {
		return elystypes.ZeroDec34(), ErrAmountTooLow
	}

	exitValue := tvl.Mul(refundedShares).Quo(elystypes.NewDec34FromInt(totalShares.Amount))

	if exitingShares.GTE(totalShares.Amount) {
		return elystypes.ZeroDec34(), errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
	}

	return exitValue, nil
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
) (exitCoins sdk.Coins, weightBalanceBonus elystypes.Dec34, err error) {
	totalShares := pool.GetTotalShares()
	if exitingShares.GTE(totalShares.Amount) {
		return sdk.Coins{}, elystypes.ZeroDec34(), errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
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
		initialWeightDistance := pool.WeightDistanceFromTarget(ctx, oracleKeeper, accountedAssets)
		tokenPrice, decimals := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
		exitValueWithoutSlippage, err := CalcExitValueWithoutSlippage(ctx, oracleKeeper, accountedPoolKeeper, pool, exitingShares, tokenOutDenom)
		if err != nil {
			return sdk.Coins{}, elystypes.ZeroDec34(), err
		}

		// Ensure tokenPrice is not zero to avoid division by zero
		if tokenPrice.IsZero() {
			return sdk.Coins{}, elystypes.ZeroDec34(), ErrAmountTooLow
		}

		oracleOutAmount := exitValueWithoutSlippage.Quo(tokenPrice.QuoInt(OneTokenUnit(decimals)))

		newAssetPools, err := pool.NewPoolAssetsAfterSwap(ctx,
			sdk.Coins{},
			sdk.Coins{sdk.NewCoin(tokenOutDenom, oracleOutAmount.ToInt())}, accountedAssets,
		)
		if err != nil {
			return sdk.Coins{}, elystypes.ZeroDec34(), err
		}
		for _, asset := range newAssetPools {
			if asset.Token.Amount.IsNegative() {
				return sdk.Coins{}, elystypes.ZeroDec34(), errors.New("out amount exceeds liquidity balance")
			}
		}

		weightDistance := pool.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
		distanceDiff := weightDistance.Sub(initialWeightDistance)

		// target weight
		targetWeightOut := GetDenomNormalizedWeight(pool.PoolAssets, tokenOutDenom)
		targetWeightIn := elystypes.OneDec34().Sub(targetWeightOut)

		// weight breaking fee as in Plasma pool
		finalWeightOut := GetDenomOracleAssetWeight(ctx, pool.PoolId, oracleKeeper, newAssetPools, tokenOutDenom)
		finalWeightIn := elystypes.OneDec34().Sub(finalWeightOut)
		initialAssetPools, err := pool.NewPoolAssetsAfterSwap(ctx,
			sdk.NewCoins(),
			sdk.NewCoins(), accountedAssets,
		)
		if err != nil {
			return sdk.Coins{}, elystypes.ZeroDec34(), err
		}
		initialWeightOut := GetDenomOracleAssetWeight(ctx, pool.PoolId, oracleKeeper, initialAssetPools, tokenOutDenom)
		initialWeightIn := elystypes.OneDec34().Sub(initialWeightOut)
		weightBreakingFee := GetWeightBreakingFee(finalWeightIn, finalWeightOut, targetWeightIn, targetWeightOut, initialWeightIn, initialWeightOut, distanceDiff, params)

		tokenOutAmount := oracleOutAmount.Mul(elystypes.OneDec34().Sub(weightBreakingFee)).ToInt()
		return sdk.Coins{sdk.NewCoin(tokenOutDenom, tokenOutAmount)}, weightBreakingFee.Neg(), nil
	}

	for _, asset := range poolLiquidity {
		// round down here, due to not wanting to over-exit
		exitAmt := shareOutRatio.MulInt(asset.Amount).TruncateInt()
		if exitAmt.LTE(sdkmath.ZeroInt()) {
			continue
		}
		if exitAmt.GTE(asset.Amount) {
			return sdk.Coins{}, elystypes.ZeroDec34(), errors.New("too many shares out")
		}
		exitedCoins = exitedCoins.Add(sdk.NewCoin(asset.Denom, exitAmt))
	}

	return exitedCoins, elystypes.ZeroDec34(), nil
}
