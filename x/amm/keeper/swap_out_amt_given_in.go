package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (k Keeper) SwapOutAmtGivenIn(
	ctx sdk.Context, poolId uint64,
	oracleKeeper types.OracleKeeper,
	snapshot types.SnapshotPool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee osmomath.BigDec,
	weightBreakingFeePerpetualFactor osmomath.BigDec,
	takersFee osmomath.BigDec,
) (tokenOut sdk.Coin, slippage osmomath.BigDec, slippageAmount osmomath.BigDec, weightBalanceBonus osmomath.BigDec, oracleOutAmount osmomath.BigDec, swapFeeFinal osmomath.BigDec, err error) {
	ammPool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), fmt.Errorf("invalid pool: %d", poolId)
	}
	params := k.GetParams(ctx)
	return ammPool.SwapOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, tokenOutDenom, swapFee, k.accountedPoolKeeper, weightBreakingFeePerpetualFactor, params, takersFee)
}
