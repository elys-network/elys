package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
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
) (sdk.Coin, elystypes.Dec34, error) {
	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), err
	}

	tokenAmountInAfterFee := elystypes.NewDec34FromInt(tokenIn.Amount).Mul(elystypes.OneDec34().SubLegacyDec(swapFee))
	poolTokenInBalance := elystypes.NewDec34FromInt(poolAssetIn.Token.Amount)
	// accounted pool balance
	acountedPoolAssetInAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetIn.Token.Denom)
	if acountedPoolAssetInAmt.IsPositive() {
		poolTokenInBalance = elystypes.NewDec34FromInt(acountedPoolAssetInAmt)
	}

	poolTokenOutBalance := elystypes.NewDec34FromInt(poolAssetOut.Token.Amount)
	// accounted pool balance
	accountedPoolAssetOutAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetOut.Token.Denom)
	if accountedPoolAssetOutAmt.IsPositive() {
		poolTokenOutBalance = elystypes.NewDec34FromInt(accountedPoolAssetOutAmt)
	}

	poolPostSwapInBalance := poolTokenInBalance.Add(tokenAmountInAfterFee)

	outWeight := elystypes.NewDec34FromInt(poolAssetOut.Weight)
	inWeight := elystypes.NewDec34FromInt(poolAssetIn.Weight)
	if p.PoolParams.UseOracle {
		_, poolAssetIn, poolAssetOut, err := snapshot.parsePoolAssets(tokensIn, tokenOutDenom)
		if err != nil {
			return sdk.Coin{}, elystypes.ZeroDec34(), err
		}
		oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, p.PoolId, oracle, []PoolAsset{poolAssetIn, poolAssetOut})
		if err != nil {
			return sdk.Coin{}, elystypes.ZeroDec34(), err
		}
		inWeight = oracleWeights[0].Weight
		outWeight = oracleWeights[1].Weight
	}

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
		return sdk.Coin{}, elystypes.ZeroDec34(), err
	}
	if tokenAmountOut.IsZero() {
		return sdk.Coin{}, elystypes.ZeroDec34(), ErrTokenOutAmountZero
	}

	rate, err := p.GetTokenARate(ctx, oracle, snapshot, tokenIn.Denom, tokenOutDenom, accountedPool)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), err
	}

	amountOutWithoutSlippage := tokenAmountInAfterFee.Mul(rate)

	// check if amountOutWithoutSlippage is zero to avoid division by zero
	if amountOutWithoutSlippage.IsZero() {
		return sdk.Coin{}, elystypes.ZeroDec34(), errorsmod.Wrapf(ErrInvalidMathApprox, "amount out without slippage must be positive")
	}

	slippage := elystypes.OneDec34().Sub(tokenAmountOut.Quo(amountOutWithoutSlippage))

	// We ignore the decimal component, as we round down the token amount out.
	tokenAmountOutInt := tokenAmountOut.ToInt()
	if !tokenAmountOutInt.IsPositive() {
		return sdk.Coin{}, elystypes.ZeroDec34(), ErrTokenOutAmountZero
	}

	return sdk.NewCoin(tokenOutDenom, tokenAmountOutInt), slippage, nil
}
