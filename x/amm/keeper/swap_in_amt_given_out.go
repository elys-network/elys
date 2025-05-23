package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
// weightBreakingFeePerpetualFactor should be 1 if perpetual is not the one calling this function
func (k Keeper) SwapInAmtGivenOut(
	ctx sdk.Context, poolId uint64, oracleKeeper types.OracleKeeper, snapshot types.SnapshotPool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee osmomath.BigDec, weightBreakingFeePerpetualFactor osmomath.BigDec, takersFee osmomath.BigDec) (
	tokenIn sdk.Coin, slippage, slippageAmount osmomath.BigDec, weightBalanceBonus osmomath.BigDec, oracleInAmount osmomath.BigDec, swapFeeFinal osmomath.BigDec, err error,
) {
	ammPool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("invalid pool: %d", poolId)
	}
	params := k.GetParams(ctx)
	return ammPool.SwapInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, swapFee, k.accountedPoolKeeper, weightBreakingFeePerpetualFactor, params, takersFee)
}
