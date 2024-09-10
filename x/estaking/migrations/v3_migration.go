package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	params := types.Params{}
	if legacyParams.StakeIncentives != nil {
		params.StakeIncentives = &types.IncentiveInfo{
			EdenAmountPerYear: legacyParams.StakeIncentives.EdenAmountPerYear,
			BlocksDistributed: legacyParams.StakeIncentives.BlocksDistributed.Int64(),
		}
	}
	params.EdenCommitVal = legacyParams.EdenCommitVal
	params.EdenbCommitVal = legacyParams.EdenbCommitVal
	params.MaxEdenRewardAprStakers = legacyParams.MaxEdenRewardAprStakers
	params.EdenBoostApr = legacyParams.EdenBoostApr
	params.DexRewardsStakers = types.DexRewardsTracker{
		NumBlocks: legacyParams.DexRewardsStakers.NumBlocks.Int64(),
		Amount:    legacyParams.DexRewardsStakers.Amount,
	}

	m.keeper.SetParams(ctx, params)
	return nil
}
