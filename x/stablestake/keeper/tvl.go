package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) TVL(ctx sdk.Context, poolId uint64) osmomath.BigDec {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return osmomath.ZeroBigDec()
	}
	netAmount := pool.GetBigDecNetAmount()
	price := k.oracleKeeper.GetDenomPrice(ctx, pool.DepositDenom)
	return price.Mul(netAmount)
}

func (k Keeper) AllTVL(ctx sdk.Context) osmomath.BigDec {
	allPools := k.GetAllPools(ctx)
	tvl := osmomath.ZeroBigDec()
	for _, pool := range allPools {
		tvl = tvl.Add(k.TVL(ctx, pool.Id))
	}
	return tvl
}
