package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) applyJoinPoolStateChange(ctx sdk.Context, pool types.Pool, joiner sdk.AccAddress, numShares sdk.Int, joinCoins sdk.Coins) error {
	err := k.bankKeeper.SendCoins(ctx, joiner, sdk.AccAddress(pool.GetAddress()), joinCoins)
	if err != nil {
		return err
	}

	err = k.MintPoolShareToAccount(ctx, pool, joiner, numShares)
	if err != nil {
		return err
	}

	err = k.SetPool(ctx, pool)
	if err != nil {
		return err
	}

	types.EmitAddLiquidityEvent(ctx, joiner, pool.GetPoolId(), joinCoins)
	k.hooks.AfterJoinPool(ctx, joiner, pool.GetPoolId(), joinCoins, numShares)
	k.RecordTotalLiquidityIncrease(ctx, joinCoins)
	return nil
}
