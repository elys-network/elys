package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// MintPoolShareToAccount attempts to mint shares of a AMM denomination to the
// specified address returning an error upon failure. Shares are minted using
// the x/amm module account.
func (k Keeper) MintPoolShareToAccount(ctx sdk.Context, pool types.Pool, addr sdk.AccAddress, amount math.Int) error {
	amt := sdk.NewCoins(sdk.NewCoin(types.GetPoolShareDenom(pool.GetPoolId()), amount))

	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, amt)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, amt)
	if err != nil {
		return err
	}

	return nil
}
