package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) PreparePools(ctx sdk.Context, tradingAsset string) (poolId uint64, ammPool ammtypes.Pool, pool types.Pool, err error) {
	poolId, err = k.GetFirstValidPool(ctx, tradingAsset)
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
		pool.InitiatePool(ctx, &ammPool)
		k.SetPool(ctx, pool)
	}

	return
}
