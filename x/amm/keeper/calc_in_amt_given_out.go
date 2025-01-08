package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcInAmtGivenOut calculates token to be provided, fee added,
// given the swapped out amount, using solveConstantFunctionInvariant.
func (k Keeper) CalcInAmtGivenOut(
	ctx sdk.Context,
	poolId uint64,
	oracle types.OracleKeeper,
	snapshot *types.Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdkmath.LegacyDec) (
	tokenIn sdk.Coin, slippage elystypes.Dec34, err error,
) {
	p, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, elystypes.ZeroDec34(), errorsmod.Wrapf(types.ErrInvalidPool, "invalid pool")
	}

	return p.CalcInAmtGivenOut(ctx, oracle, snapshot, tokensOut, tokenInDenom, swapFee, k.accountedPoolKeeper)
}
