package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {
	// reset params
	params := m.keeper.GetParams(ctx)
	params.TakerFees = math.LegacyZeroDec()
	m.keeper.SetParams(ctx, params)

	return nil
}
