package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	params := types.Params{}
	if legacyParams.LpIncentives != nil {
		params.LpIncentives = &types.IncentiveInfo{
			EdenAmountPerYear: legacyParams.LpIncentives.EdenAmountPerYear,
			BlocksDistributed: legacyParams.LpIncentives.BlocksDistributed.Int64(),
		}
	}
	params.MaxEdenRewardAprLps = legacyParams.MaxEdenRewardAprLps
	params.RewardPortionForLps = legacyParams.RewardPortionForLps
	params.RewardPortionForStakers = legacyParams.RewardPortionForStakers
	params.SupportedRewardDenoms = legacyParams.SupportedRewardDenoms
	params.ProtocolRevenueAddress = legacyParams.ProtocolRevenueAddress

	m.keeper.SetParams(ctx, params)

	// Migrate pool infos
	pools := m.keeper.GetAllLegacyPoolInfos(ctx)

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

		m.keeper.RemoveLegacyPoolInfo(ctx, pool.PoolId)
		m.keeper.SetPoolInfo(ctx, newPool)
	}
	return nil
}
