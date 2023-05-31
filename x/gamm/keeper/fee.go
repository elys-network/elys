package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	poolmanagertypes "github.com/elys-network/elys/x/poolmanager/types"
)

func portionCoins(coins sdk.Coins, portion sdk.Dec) sdk.Coins {
	portionCoins := sdk.Coins{}
	for _, coin := range coins {
		portionCoins = portionCoins.Add(sdk.NewCoin(
			coin.Denom, coin.Amount.ToDec().Mul(sdk.OneDec().Sub(portion)).RoundInt(),
		))
	}
	return portionCoins
}

func (k Keeper) OnCollectFee(ctx sdk.Context, pool poolmanagertypes.PoolI, fee sdk.Coins) {
	// TODO: transfer 60% to LPs through commitment module
	// TODO: transfer 30% to stakers through incentive module
	// keep remaining 10% on fees treasury
}
