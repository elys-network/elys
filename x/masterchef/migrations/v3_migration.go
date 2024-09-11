package migrations

/*
func (m Migrator) V3Migration(ctx sdk.Context) error {
	params := m.keeper.GetLegacyParams(ctx)
	m.keeper.SetParams(ctx, params)
	m.keeper.DeleteLegacyParams(ctx)

	//m.keeper.MigrateFromV2UserRewardInfos(ctx)

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

	legacyPoolRewardAccums := m.keeper.GetAllLegacyPoolRewardsAccum(ctx)
	for _, legacyPoolRewardAccum := range legacyPoolRewardAccums {
		m.keeper.SetPoolRewardsAccum(ctx, legacyPoolRewardAccum)
		m.keeper.DeleteLegacyPoolRewardsAccum(ctx, legacyPoolRewardAccum)
	}

	return nil
}*/
