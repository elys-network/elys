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
		tokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
		exitValueWithoutSlippage, err := CalcExitValueWithoutSlippage(ctx, oracleKeeper, accountedPoolKeeper, pool, exitingShares, tokenOutDenom)
		if err != nil {
			return sdk.Coins{}, elystypes.ZeroDec34(), err
		}

		// Ensure tokenPrice is not zero to avoid division by zero
		if tokenPrice.IsZero() {
			return sdk.Coins{}, elystypes.ZeroDec34(), ErrAmountTooLow
		}

		oracleOutAmount := exitValueWithoutSlippage.Quo(tokenPrice)

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

		weightBalanceBonus, weightBreakingFee, _ := pool.CalculateWeightFees(ctx, oracleKeeper, accountedAssets, newAssetPools, tokenOutDenom, params, sdkmath.LegacyOneDec())
		// apply percentage to fees, consider improvement or reduction of other token
		// Other denom weight ratio to reduce the weight breaking fees
		initialWeightOut := GetDenomOracleAssetWeight(ctx, pool.PoolId, oracleKeeper, accountedAssets, tokenOutDenom)
		initialWeightIn := elystypes.OneDec34().Sub(initialWeightOut)
		weightBreakingFee = weightBreakingFee.Mul(initialWeightIn)

		tokenOutAmount := oracleOutAmount.Mul(elystypes.OneDec34().Sub(weightBreakingFee)).ToInt()
		return sdk.Coins{sdk.NewCoin(tokenOutDenom, tokenOutAmount)}, weightBalanceBonus, nil
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
