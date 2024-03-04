package keeper

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (k Keeper) SwapOutAmtGivenIn(
	ctx sdk.Context, poolId uint64,
	oracleKeeper types.OracleKeeper,
	snapshot *types.Pool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee sdk.Dec,
) (tokenOut sdk.Coin, slippage sdk.Dec, slippageAmount sdk.Dec, weightBalanceBonus sdk.Dec, err error) {
	ammPool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec(), fmt.Errorf("invalid pool: %d", poolId)
	}

	return ammPool.SwapOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, tokenOutDenom, swapFee, k.accountedPoolKeeper)
}
