package types

import (
	"errors"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CalcExitPool returns how many tokens should come out, when exiting k LP shares against a "standard" CFMM
func CalcExitPool(ctx sdk.Context, oracleKeeper OracleKeeper, pool Pool, exitingShares sdk.Int, tokenOutDenom string) (sdk.Coins, error) {
	totalShares := pool.GetTotalShares()
	if exitingShares.GTE(totalShares.Amount) {
		return sdk.Coins{}, sdkerrors.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
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
		tvl, err := pool.TVL(ctx, oracleKeeper)
		if err != nil {
			return sdk.Coins{}, err
		}

		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
		oracleOutAmount := tvl.Mul(refundedShares).Quo(refundedShares).Quo(tokenPrice)

		newAssetPools, err := pool.NewPoolAssetsAfterSwap(
			sdk.Coins{},
			sdk.Coins{sdk.NewCoin(tokenOutDenom, oracleOutAmount.RoundInt())},
		)
		if err != nil {
			return sdk.Coins{}, err
		}
		for _, asset := range newAssetPools {
			if asset.Token.Amount.IsNegative() {
				return sdk.Coins{}, fmt.Errorf("out amount exceeds liquidity balance")
			}
		}

		weightDistance := pool.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
		distanceDiff := weightDistance.Sub(initialWeightDistance)
		weightBreakingFee := sdk.ZeroDec()
		if distanceDiff.IsPositive() {
			weightBreakingFee = pool.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff)
		}

		tokenOutAmount := oracleOutAmount.Mul(sdk.OneDec().Sub(weightBreakingFee)).RoundInt()
		return sdk.Coins{sdk.NewCoin(tokenOutDenom, tokenOutAmount)}, nil
	}

	for _, asset := range poolLiquidity {
		// round down here, due to not wanting to over-exit
		exitAmt := shareOutRatio.MulInt(asset.Amount).TruncateInt()
		if exitAmt.LTE(sdk.ZeroInt()) {
			continue
		}
		if exitAmt.GTE(asset.Amount) {
			return sdk.Coins{}, errors.New("too many shares out")
		}
		exitedCoins = exitedCoins.Add(sdk.NewCoin(asset.Denom, exitAmt))
	}

	return exitedCoins, nil
}
