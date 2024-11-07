package keeper

import (
	"cosmossdk.io/math"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
// weightBreakingFeePerpetualFactor should be 1 if perpetual is not the one calling this function
func (k Keeper) SwapInAmtGivenOut(
	ctx sdk.Context, poolId uint64, oracleKeeper types.OracleKeeper, snapshot *types.Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec, weightBreakingFeePerpetualFactor math.LegacyDec) (
	tokenIn sdk.Coin, slippage, slippageAmount sdk.Dec, weightBalanceBonus sdk.Dec, err error,
) {
	ammPool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), fmt.Errorf("invalid pool: %d", poolId)
	}

	return ammPool.SwapInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, swapFee, k.accountedPoolKeeper, weightBreakingFeePerpetualFactor)
}
