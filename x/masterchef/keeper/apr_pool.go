package keeper

import (
	sdkmath "cosmossdk.io/math"
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
				PoolId:     poolId,
				UsdcDexApr: sdkmath.LegacyZeroDec(),
				UsdcGasApr: sdkmath.LegacyZeroDec(),
				EdenApr:    sdkmath.LegacyZeroDec(),
				TotalApr:   sdkmath.LegacyZeroDec(),
			})
			continue
		}
		data = append(data, types.PoolApr{
			PoolId:     poolId,
			UsdcDexApr: poolInfo.DexApr,
			UsdcGasApr: poolInfo.GasApr,
			EdenApr:    poolInfo.EdenApr,
			TotalApr:   poolInfo.DexApr.Add(poolInfo.GasApr).Add(poolInfo.EdenApr),
		})
	}

	return data
}
