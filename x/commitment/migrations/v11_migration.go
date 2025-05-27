package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	// May 29th 7pm UTC approx based 3.8s block time
	params.StartRewardProgramClaimHeight = 3_908_564
	// After 31 days of the start height
	params.EndRewardProgramClaimHeight = 4_613_406
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
