package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) CalculatePoolAprs(ctx sdk.Context, ids []uint64) []types.PoolApr {
	if len(ids) == 0 {
		pools := k.GetAllPoolInfos(ctx)
		for _, pool := range pools {
			ids = append(ids, pool.PoolId)
		}
	}

	data := []types.PoolApr{}
	for _, poolId := range ids {
		poolInfo, found := k.GetPoolInfo(ctx, poolId)
		if !found {
			data = append(data, types.PoolApr{
				PoolId:   poolId,
				UsdcApr:  sdk.ZeroDec(),
				EdenApr:  sdk.ZeroDec(),
				TotalApr: sdk.ZeroDec(),
			})
			continue
		}
		data = append(data, types.PoolApr{
			PoolId:   poolId,
			UsdcApr:  poolInfo.DexApr,
			EdenApr:  poolInfo.EdenApr,
			TotalApr: poolInfo.DexApr.Add(poolInfo.EdenApr),
		})
	}

	return data
}
