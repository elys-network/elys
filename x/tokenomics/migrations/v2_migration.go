package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/tokenomics/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	// reset genesis inflation param
	inflation := types.GenesisInflation{
		Authority:             "elys10d07y265gmmuvt4z0w9aw880jnsr700j6z2zm3",
		SeedVesting:           6000000000000,
		StrategicSalesVesting: 180000000000000,
		Inflation: &types.InflationEntry{
			CommunityFund:     9750000000000,
			IcsStakingRewards: 0,
			LmRewards:         0,
			StrategicReserve:  18750000000000,
			TeamTokensVested:  0,
		},
	}

	m.keeper.SetGenesisInflation(ctx, inflation)
	return nil
}
