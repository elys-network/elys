package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CalcOutAmtGivenIn calculates tokens to be swapped out given the provided
// amount and fee deducted, using solveConstantFunctionInvariant.
func (p Pool) CalcOutAmtGivenIn(
	ctx sdk.Context,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee sdk.Dec,
) (sdk.Coin, error) {
	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, err
	}

	tokenAmountInAfterFee := sdk.NewDecFromInt(tokenIn.Amount).Mul(sdk.OneDec().Sub(swapFee))
	poolTokenInBalance := sdk.NewDecFromInt(poolAssetIn.Token.Amount)
	poolPostSwapInBalance := poolTokenInBalance.Add(tokenAmountInAfterFee)

	// deduct swapfee on the tokensIn
	// delta balanceOut is positive(tokens inside the pool decreases)
	tokenAmountOut := solveConstantFunctionInvariant(
		poolTokenInBalance,
		poolPostSwapInBalance,
		sdk.NewDecFromInt(poolAssetIn.Weight),
		sdk.NewDecFromInt(poolAssetOut.Token.Amount),
		sdk.NewDecFromInt(poolAssetOut.Weight),
	)

	// We ignore the decimal component, as we round down the token amount out.
	tokenAmountOutInt := tokenAmountOut.TruncateInt()
	if !tokenAmountOutInt.IsPositive() {
		return sdk.Coin{}, sdkerrors.Wrapf(ErrInvalidMathApprox, "token amount must be positive")
	}

	return sdk.NewCoin(tokenOutDenom, tokenAmountOutInt), nil
}
