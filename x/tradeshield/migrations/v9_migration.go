package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V9Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.Tolerance = math.LegacyMustNewDecFromStr("0.05")
	m.keeper.SetParams(ctx, &params)

	return nil
}
