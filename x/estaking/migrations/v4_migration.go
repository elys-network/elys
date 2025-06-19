package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochsmoduletypes "github.com/elys-network/elys/v6/x/epochs/types"
	"github.com/elys-network/elys/v6/x/estaking/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	params := types.Params{}
	if legacyParams.StakeIncentives != nil {
		params.StakeIncentives = &types.IncentiveInfo{
			EdenAmountPerYear: legacyParams.StakeIncentives.EdenAmountPerYear,
			BlocksDistributed: legacyParams.StakeIncentives.BlocksDistributed,
		}
	}
	params.EdenCommitVal = legacyParams.EdenCommitVal
	params.EdenbCommitVal = legacyParams.EdenbCommitVal
	params.MaxEdenRewardAprStakers = legacyParams.MaxEdenRewardAprStakers
	params.EdenBoostApr = legacyParams.EdenBoostApr
	params.ProviderVestingEpochIdentifier = epochsmoduletypes.TenDaysEpochID
	params.ProviderStakingRewardsPortion = math.LegacyMustNewDecFromStr("0.25")

	m.keeper.SetParams(ctx, params)
	return nil
}
