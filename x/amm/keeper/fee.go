package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func portionCoins(coins sdk.Coins, portion sdk.Dec) sdk.Coins {
	portionCoins := sdk.Coins{}
	for _, coin := range coins {
		portionAmount := sdk.NewDecFromInt(coin.Amount).Mul(portion).RoundInt()
		portionCoins = portionCoins.Add(sdk.NewCoin(
			coin.Denom, portionAmount,
		))
	}
	return portionCoins
}

func (k Keeper) OnCollectFee(ctx sdk.Context, pool types.Pool, fee sdk.Coins) {
	// TODO: transfer 60% to LPs through commitment module
	// TODO: transfer 30% to stakers through incentive module
	// keep remaining 10% on fees treasury
	// TODO: swap fees to pool FeeDenom (this is normally USDC)
}
