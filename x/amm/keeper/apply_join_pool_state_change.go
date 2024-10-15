package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) applyJoinPoolStateChange(ctx sdk.Context, pool types.Pool, joiner sdk.AccAddress, numShares math.Int, joinCoins sdk.Coins) error {
	if err := k.bankKeeper.SendCoins(ctx, joiner, sdk.MustAccAddressFromBech32(pool.GetAddress()), joinCoins); err != nil {
		return err
	}

	if err := k.MintPoolShareToAccount(ctx, pool, joiner, numShares); err != nil {
		return err
	}

	if err := k.SetPool(ctx, pool); err != nil {
		return err
	}

	types.EmitAddLiquidityEvent(ctx, joiner, pool.GetPoolId(), joinCoins)
	if k.hooks != nil {
		return k.hooks.AfterJoinPool(ctx, joiner, pool, joinCoins, numShares)
	}
	k.RecordTotalLiquidityIncrease(ctx, joinCoins)
	return nil
}
