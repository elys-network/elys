package keeper

import (
	"fmt"

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

func (k Keeper) GetTotalAndPerDenomTVL(ctx sdk.Context) (totalTVL osmomath.BigDec, denomTVL sdk.Coins, err error) {
	allPools := k.GetAllPools(ctx)
	totalTVL = osmomath.ZeroBigDec()
	denomTVL = sdk.Coins{}
	for _, pool := range allPools {
		poolTVL := k.TVL(ctx, pool.Id)
		totalTVL = totalTVL.Add(poolTVL)
		assetInfo, found := k.oracleKeeper.GetAssetInfo(ctx, pool.DepositDenom)
		if !found {
			return osmomath.BigDec{}, sdk.Coins{}, fmt.Errorf("asset info %s not found", pool.DepositDenom)
		}
		denomTVL = denomTVL.Add(sdk.NewCoin(assetInfo.Display, poolTVL.Dec().TruncateInt()))
	}
	return totalTVL, denomTVL, nil
}
