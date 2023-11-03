package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (k Keeper) SwapInAmtGivenOut(
	ctx sdk.Context, poolId uint64, oracleKeeper types.OracleKeeper, snapshot *types.Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (
	tokenIn sdk.Coin, slippageAmount sdk.Dec, weightBalanceBonus sdk.Dec, err error,
) {
	ammPool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), fmt.Errorf("invalid pool: %d", poolId)
	}

	return ammPool.SwapInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, swapFee, k.accountedPoolKeeper)
}
