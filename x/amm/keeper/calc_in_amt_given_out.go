package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CalcInAmtGivenOut calculates token to be provided, fee added,
// given the swapped out amount, using solveConstantFunctionInvariant.
func (k Keeper) CalcInAmtGivenOut(
	ctx sdk.Context,
	poolId uint64,
	oracle types.OracleKeeper,
	snapshot types.SnapshotPool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee osmomath.BigDec) (
	tokenIn sdk.Coin, slippage osmomath.BigDec, err error,
) {
	p, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, osmomath.ZeroBigDec(), errorsmod.Wrapf(types.ErrInvalidPool, "invalid pool")
	}

	return p.CalcInAmtGivenOut(ctx, oracle, snapshot, tokensOut, tokenInDenom, swapFee)
}
