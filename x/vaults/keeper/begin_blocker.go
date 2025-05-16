package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/vaults/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// Traverse all vaults and deduct management fee from all coins and send it to the vault's manager and protocol revenue address
	vaults := k.GetAllVaults(ctx)
	for _, vault := range vaults {
		coins := k.bk.GetAllBalances(ctx, types.NewVaultAddress(vault.Id))
		for _, coin := range coins {
			// TODO: Fix this, deduct fee on per year basis
			coin.Amount = coin.Amount.ToLegacyDec().Mul(vault.ManagementFee).RoundInt()
		}
		// send coins to protocol revenue address and manager address
		err := k.bk.SendCoins(ctx, types.NewVaultAddress(vault.Id), sdk.MustAccAddressFromBech32(vault.Manager), coins)
		if err != nil {
			// log error
			k.Logger().Error("error sending coins to vault manager", "error", err)
		}
	}
}
