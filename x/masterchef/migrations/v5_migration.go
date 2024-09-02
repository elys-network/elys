package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllPoolInfos(ctx)

	ctx.Logger().Info("Migration: Adding enable eden rewards field")

	for _, pool := range pools {

		ctx.Logger().Debug("pool", pool)

		newPool := types.PoolInfo{
			PoolId:               pool.PoolId,
			RewardWallet:         pool.RewardWallet,
			Multiplier:           pool.Multiplier,
			EdenApr:              pool.EdenApr,
			DexApr:               pool.DexApr,
			GasApr:               pool.GasApr,
			ExternalIncentiveApr: pool.ExternalIncentiveApr,
			ExternalRewardDenoms: pool.ExternalRewardDenoms,
			EnableEdenRewards:    true,
		}

		m.keeper.RemovePoolInfo(ctx, pool.PoolId)
		m.keeper.SetPoolInfo(ctx, newPool)
	}
	return nil
}
