package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	// reset params
	legacy := m.keeper.GetLegacyParams(ctx)
	params := types.NewParams(
		legacy.MinCommissionRate, // min commission 0.05
		legacy.MaxVotingPower,    // max voting power 0.66
		legacy.MinSelfDelegation, // min self delegation
		legacy.BrokerAddress,     // broker address
		6307200,                  // total blocks per year
	)
	m.keeper.SetParams(ctx, params)

	return nil
}
