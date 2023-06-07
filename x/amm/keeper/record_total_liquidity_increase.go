package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) RecordTotalLiquidityIncrease(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		_, found := k.GetDenomLiquidity(ctx, coin.Denom)
		if !found {
			k.SetDenomLiquidity(ctx, types.DenomLiquidity{Denom: coin.Denom, Liquidity: sdk.ZeroInt()})
		}
		k.IncreaseDenomLiquidity(ctx, coin.Denom, coin.Amount)
	}
}
