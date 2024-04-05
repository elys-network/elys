package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)

	totalCommitted := sdk.Coins{}
	legacyCommitments := m.keeper.GetAllLegacyCommitments(ctx)
	for _, legacy := range legacyCommitments {

		vestingTokens := []*types.VestingTokens{}
		for _, vt := range legacy.VestingTokens {
			blockMultiplier := int64(17280) // "day"
			if vt.EpochIdentifier == "tenseconds" {
				blockMultiplier = int64(2)
			} else {
				blockMultiplier = int64(720) // "hour"
			}
			vestingTokens = append(vestingTokens, &types.VestingTokens{
				Denom:                vt.Denom,
				TotalAmount:          vt.UnvestedAmount,
				ClaimedAmount:        math.ZeroInt(),
				NumBlocks:            vt.NumEpochs * blockMultiplier,
				StartBlock:           ctx.BlockHeight(),
				VestStartedTimestamp: vt.VestStartedTimestamp,
			})
		}
		m.keeper.SetCommitments(ctx, types.Commitments{
			Creator:                 legacy.Creator,
			CommittedTokens:         legacy.CommittedTokens,
			RewardsUnclaimed:        legacy.RewardsUnclaimed,
			Claimed:                 legacy.Claimed,
			VestingTokens:           vestingTokens,
			RewardsByElysUnclaimed:  legacy.RewardsByElysUnclaimed,
			RewardsByEdenUnclaimed:  legacy.RewardsByEdenUnclaimed,
			RewardsByEdenbUnclaimed: legacy.RewardsByEdenbUnclaimed,
			RewardsByUsdcUnclaimed:  legacy.RewardsByUsdcUnclaimed,
		})
		for _, token := range legacy.CommittedTokens {
			totalCommitted = totalCommitted.Add(sdk.NewCoin(token.Denom, token.Amount))
		}
	}

	vestingInfos := []*types.VestingInfo{}
	for _, legacy := range legacyParams.VestingInfos {
		vestingInfos = append(vestingInfos, &types.VestingInfo{
			BaseDenom:      legacy.BaseDenom,
			VestingDenom:   legacy.VestingDenom,
			NumBlocks:      legacy.NumEpochs * 17280,
			VestNowFactor:  legacy.VestNowFactor,
			NumMaxVestings: legacy.NumMaxVestings,
		})
	}

	m.keeper.SetParams(ctx, types.Params{
		VestingInfos:   vestingInfos,
		TotalCommitted: totalCommitted,
	})
	return nil
}
