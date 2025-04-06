package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) TVL(ctx sdk.Context, oracleKeeper types.OracleKeeper, poolId uint64) osmomath.BigDec {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return osmomath.ZeroBigDec()
	}
	totalDeposit := pool.GetBigDecTotalValue()
	price := oracleKeeper.GetDenomPrice(ctx, pool.DepositDenom)
	return price.Mul(totalDeposit)
}

func (k Keeper) AllTVL(ctx sdk.Context, oracleKeeper types.OracleKeeper) osmomath.BigDec {
	allPools := k.GetAllPools(ctx)
	tvl := osmomath.ZeroBigDec()
	for _, pool := range allPools {
		tvl = tvl.Add(k.TVL(ctx, oracleKeeper, pool.Id))
	}
	return tvl
}
