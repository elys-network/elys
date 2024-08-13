package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	legacyUserRewardInfos := m.keeper.GetAllLegacyUserRewardInfos(ctx)

	for _, legacyUserRewardInfo := range legacyUserRewardInfos {
		m.keeper.SetUserRewardInfo(ctx, legacyUserRewardInfo)
		m.keeper.DeleteLegacyUserRewardInfo(ctx, legacyUserRewardInfo.User, legacyUserRewardInfo.PoolId, legacyUserRewardInfo.RewardDenom)
	}

	legacyPoolInfos := m.keeper.GetAllLegacyPoolInfos(ctx)
	for _, legacyPoolInfo := range legacyPoolInfos {
		m.keeper.SetPoolInfo(ctx, legacyPoolInfo)
		m.keeper.RemoveLegacyPoolInfo(ctx, legacyPoolInfo.PoolId)
	}

	externalIncentiveIndex := m.keeper.GetLegacyExternalIncentiveIndex(ctx)
	m.keeper.SetExternalIncentiveIndex(ctx, externalIncentiveIndex)
	m.keeper.RemoveLegacyExternalIncentiveIndex(ctx)

	externalIncentives := m.keeper.GetAllLegacyExternalIncentives(ctx)
	for _, externalIncentive := range externalIncentives {
		m.keeper.SetExternalIncentive(ctx, externalIncentive)
		m.keeper.DeleteLegacyExternalIncentive(ctx, externalIncentive.Id)
	}

	legacyPoolRewardInfos := m.keeper.GetAllLegacyPoolRewardInfos(ctx)
	for _, legacyPoolRewardInfo := range legacyPoolRewardInfos {
		m.keeper.SetPoolRewardInfo(ctx, legacyPoolRewardInfo)
		m.keeper.DeleteLegacyPoolRewardInfo(ctx, legacyPoolRewardInfo.PoolId, legacyPoolRewardInfo.RewardDenom)
	}

	return nil
}
