package keeper

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
)

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
// weightBreakingFeePerpetualFactor should be 1 if perpetual is not the one calling this function
func (k Keeper) SwapInAmtGivenOut(
	ctx sdk.Context, poolId uint64, oracleKeeper types.OracleKeeper, snapshot *types.Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee math.LegacyDec, weightBreakingFeePerpetualFactor math.LegacyDec) (
	tokenIn sdk.Coin, slippage, slippageAmount, weightBalanceBonus, oracleInAmount elystypes.Dec34, swapFeeFinal math.LegacyDec, err error,
) {
	ammPool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), math.LegacyZeroDec(), fmt.Errorf("invalid pool: %d", poolId)
	}
	params := k.GetParams(ctx)
	return ammPool.SwapInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, swapFee, k.accountedPoolKeeper, weightBreakingFeePerpetualFactor, params)
}
