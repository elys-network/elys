package types

import (
	errorsmod "cosmossdk.io/errors"
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
	swapFee sdk.Dec,
	accountedPool AccountedPoolKeeper,
) (sdk.Coin, sdk.Dec, error) {
	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	tokenAmountInAfterFee := sdk.NewDecFromInt(tokenIn.Amount).Mul(sdk.OneDec().Sub(swapFee))
	poolTokenInBalance := sdk.NewDecFromInt(poolAssetIn.Token.Amount)
	// accounted pool balance
	acountedPoolAssetInAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetIn.Token.Denom)
	if acountedPoolAssetInAmt.IsPositive() {
		poolTokenInBalance = sdk.NewDecFromInt(acountedPoolAssetInAmt)
	}

	poolTokenOutBalance := sdk.NewDecFromInt(poolAssetOut.Token.Amount)
	// accounted pool balance
	accountedPoolAssetOutAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetOut.Token.Denom)
	if accountedPoolAssetOutAmt.IsPositive() {
		poolTokenOutBalance = sdk.NewDecFromInt(accountedPoolAssetOutAmt)
	}

	poolPostSwapInBalance := poolTokenInBalance.Add(tokenAmountInAfterFee)

	outWeight := sdk.NewDecFromInt(poolAssetOut.Weight)
	inWeight := sdk.NewDecFromInt(poolAssetIn.Weight)
	if p.PoolParams.UseOracle {
		_, poolAssetIn, poolAssetOut, err := snapshot.parsePoolAssets(tokensIn, tokenOutDenom)
		if err != nil {
			return sdk.Coin{}, sdk.ZeroDec(), err
		}
		oracleWeights, err := OraclePoolNormalizedWeights(ctx, oracle, []PoolAsset{poolAssetIn, poolAssetOut})
		if err != nil {
			return sdk.Coin{}, sdk.ZeroDec(), err
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
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	rate, err := p.GetTokenARate(ctx, oracle, snapshot, tokenIn.Denom, tokenOutDenom, accountedPool)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	amountOutWithoutSlippage := tokenAmountInAfterFee.Mul(rate)
	slippage := sdk.OneDec().Sub(tokenAmountOut.Quo(amountOutWithoutSlippage))

	// We ignore the decimal component, as we round down the token amount out.
	tokenAmountOutInt := tokenAmountOut.TruncateInt()
	if !tokenAmountOutInt.IsPositive() {
		return sdk.Coin{}, sdk.ZeroDec(), errorsmod.Wrapf(ErrInvalidMathApprox, "token amount must be positive")
	}

	return sdk.NewCoin(tokenOutDenom, tokenAmountOutInt), slippage, nil
}
