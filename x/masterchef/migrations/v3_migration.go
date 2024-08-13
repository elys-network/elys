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

	return nil
}
