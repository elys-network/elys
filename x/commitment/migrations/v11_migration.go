package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	// TODO: Set to actual height
	params.StartRewardProgramClaimHeight = 1000000
	params.EndRewardProgramClaimHeight = 1000000
	m.keeper.SetParams(ctx, params)

	for _, rewardProgram := range RewardProgram {
		m.keeper.SetRewardProgram(ctx, rewardProgram)
	}

	if ctx.ChainID() == "elys-testnet-1" {
		for _, rewardProgram := range RewardProgramTestnet {
			m.keeper.SetRewardProgram(ctx, rewardProgram)
		}
	}

	return nil
}
