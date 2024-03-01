package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcInAmtGivenOut calculates token to be provided, fee added,
// given the swapped out amount, using solveConstantFunctionInvariant.
func (k Keeper) CalcInAmtGivenOut(
	ctx sdk.Context,
	poolId uint64,
	oracle types.OracleKeeper,
	snapshot *types.Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (
	tokenIn sdk.Coin, slippage sdk.Dec, err error,
) {
	p, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, sdk.ZeroDec(), errorsmod.Wrapf(types.ErrInvalidPool, "invalid pool")
	}

	return p.CalcInAmtGivenOut(ctx, oracle, snapshot, tokensOut, tokenInDenom, swapFee, k.accountedPoolKeeper)
}
