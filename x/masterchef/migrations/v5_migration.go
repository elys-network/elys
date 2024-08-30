package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	params := types.Params{}
	params.LpIncentives = &types.IncentiveInfo{
		EdenAmountPerYear: legacyParams.LpIncentives.EdenAmountPerYear,
		BlocksDistributed: legacyParams.LpIncentives.BlocksDistributed.Int64(),
	}
	params.MaxEdenRewardAprLps = legacyParams.MaxEdenRewardAprLps
	params.RewardPortionForLps = legacyParams.RewardPortionForLps
	params.RewardPortionForStakers = legacyParams.RewardPortionForStakers
	params.SupportedRewardDenoms = legacyParams.SupportedRewardDenoms
	params.ProtocolRevenueAddress = legacyParams.ProtocolRevenueAddress

	m.keeper.SetParams(ctx, params)
	return nil
}
