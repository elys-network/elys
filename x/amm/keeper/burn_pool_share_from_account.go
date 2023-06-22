package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// BurnPoolShareFromAccount burns `amount` of the given pool's shares held by `addr`.
func (k Keeper) BurnPoolShareFromAccount(ctx sdk.Context, pool types.Pool, addr sdk.AccAddress, amount sdk.Int) error {
	coin := sdk.NewCoin(types.GetPoolShareDenom(pool.GetPoolId()), amount)
	coins := sdk.NewCoins(coin)

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins); err != nil {
		return err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	return nil
}
