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

	k.SetPool(ctx, pool)

	types.EmitAddLiquidityEvent(ctx, joiner, pool.GetPoolId(), joinCoins)
	if k.hooks != nil {
		err := k.hooks.AfterJoinPool(ctx, joiner, pool, joinCoins, numShares)
		if err != nil {
			return err
		}
	}
	return nil
}
