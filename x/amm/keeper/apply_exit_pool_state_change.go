package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) applyExitPoolStateChange(ctx sdk.Context, pool types.Pool, exiter sdk.AccAddress, numShares sdk.Int, exitCoins sdk.Coins) error {
	err := k.bankKeeper.SendCoins(ctx, sdk.AccAddress(pool.GetAddress()), exiter, exitCoins)
	if err != nil {
		return err
	}

	err = k.BurnPoolShareFromAccount(ctx, pool, exiter, numShares)
	if err != nil {
		return err
	}

	err = k.SetPool(ctx, pool)
	if err != nil {
		return err
	}

	types.EmitRemoveLiquidityEvent(ctx, exiter, pool.GetPoolId(), exitCoins)
	k.hooks.AfterExitPool(ctx, exiter, pool.GetPoolId(), numShares, exitCoins)
	k.RecordTotalLiquidityDecrease(ctx, exitCoins)
	return nil
}
