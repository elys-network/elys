package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) TVL(ctx sdk.Context, oracleKeeper types.OracleKeeper, poolId uint64) math.LegacyDec {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return math.LegacyZeroDec()
	}
	totalDeposit := pool.TotalValue
	price := oracleKeeper.GetAssetPriceFromDenom(ctx, pool.DepositDenom)
	return price.MulInt(totalDeposit)
}

func (k Keeper) AllTVL(ctx sdk.Context, oracleKeeper types.OracleKeeper) math.LegacyDec {
	allPools := k.GetAllPools(ctx)
	tvl := math.LegacyZeroDec()
	for _, pool := range allPools {
		tvl = tvl.Add(k.TVL(ctx, oracleKeeper, pool.Id))
	}
	return tvl
}
