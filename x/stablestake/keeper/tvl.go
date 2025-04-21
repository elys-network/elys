package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) TVL(ctx sdk.Context, poolId uint64) math.LegacyDec {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return math.LegacyZeroDec()
	}
	netAmount := pool.NetAmount
	price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, pool.DepositDenom)
	return price.MulInt(netAmount)
}

func (k Keeper) AllTVL(ctx sdk.Context) math.LegacyDec {
	allPools := k.GetAllPools(ctx)
	tvl := math.LegacyZeroDec()
	for _, pool := range allPools {
		tvl = tvl.Add(k.TVL(ctx, pool.Id))
	}
	return tvl
}
