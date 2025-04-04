package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) TVL(ctx sdk.Context, oracleKeeper types.OracleKeeper, poolId uint64) elystypes.Dec34 {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return elystypes.ZeroDec34()
	}
	totalDeposit := pool.TotalValue
	price, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, pool.DepositDenom)
	return price.MulInt(totalDeposit)
}

func (k Keeper) AllTVL(ctx sdk.Context, oracleKeeper types.OracleKeeper) elystypes.Dec34 {
	allPools := k.GetAllPools(ctx)
	tvl := elystypes.ZeroDec34()
	for _, pool := range allPools {
		tvl = tvl.Add(k.TVL(ctx, oracleKeeper, pool.Id))
	}
	return tvl
}
