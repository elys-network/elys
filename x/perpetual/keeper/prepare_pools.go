package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// PreparePools creates accounted pools
func (k Keeper) PreparePools(ctx sdk.Context, collateralAsset, tradingAsset string) (poolId uint64, ammPool ammtypes.Pool, pool types.Pool, err error) {
	poolId, err = k.GetBestPool(ctx, collateralAsset, tradingAsset)
	if err != nil {
		return
	}

	ammPool, err = k.GetAmmPool(ctx, poolId, tradingAsset)
	if err != nil {
		return
	}

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		pool = types.NewPool(poolId)
		err = pool.InitiatePool(ctx, &ammPool)
		if err != nil {
			return
		}
		k.SetPool(ctx, pool)
	}

	return
}
