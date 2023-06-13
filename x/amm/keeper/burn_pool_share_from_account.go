package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// BurnPoolShareFromAccount burns `amount` of the given pools shares held by `addr`.
func (k Keeper) BurnPoolShareFromAccount(ctx sdk.Context, pool types.Pool, addr sdk.AccAddress, amount sdk.Int) error {
	amt := sdk.Coins{
		sdk.NewCoin(types.GetPoolShareDenom(pool.GetPoolId()), amount),
	}

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, amt)
	if err != nil {
		return err
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, amt)
	if err != nil {
		return err
	}

	return nil
}
