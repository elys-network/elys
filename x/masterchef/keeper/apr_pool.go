package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) CalculatePoolAprs(ctx sdk.Context, ids []uint64) []types.PoolApr {
	if len(ids) == 0 {
		pools := k.amm.GetAllPool(ctx)
		for _, pool := range pools {
			ids = append(ids, pool.PoolId)
		}
	}

	data := []types.PoolApr{}
	for _, poolId := range ids {
		poolInfo, found := k.GetPool(ctx, poolId)
		if !found {
			data = append(data, types.PoolApr{
				PoolId: poolId,
				Apr:    sdk.ZeroDec(),
			})
			continue
		}
		data = append(data, types.PoolApr{
			PoolId: poolId,
			Apr:    poolInfo.DexApr.Add(poolInfo.EdenApr),
		})
	}

	return data
}
