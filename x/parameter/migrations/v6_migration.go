package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {
	// reset params
	params := m.keeper.GetParams(ctx)
	params.TakerFees = math.LegacyZeroDec()
	// TODO: Add revenue address
	params.ProtocolRevenueAddress = ""
	m.keeper.SetParams(ctx, params)

	return nil
}
