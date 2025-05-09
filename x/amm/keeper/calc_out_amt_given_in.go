package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CalcOutAmtGivenIn calculates tokens to be swapped out given the provided
// amount and fee deducted, using solveConstantFunctionInvariant.
func (k Keeper) CalcOutAmtGivenIn(
	ctx sdk.Context,
	poolId uint64,
	oracle types.OracleKeeper,
	snapshot types.SnapshotPool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee osmomath.BigDec,
) (sdk.Coin, osmomath.BigDec, error) {
	p, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, osmomath.ZeroBigDec(), errorsmod.Wrapf(types.ErrInvalidPool, "invalid pool")
	}

	return p.CalcOutAmtGivenIn(ctx, oracle, snapshot, tokensIn, tokenOutDenom, swapFee)
}
