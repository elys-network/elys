package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V12Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	legacyParams.MinSlippage = math.LegacyMustNewDecFromStr("0.00075")

	m.keeper.SetParams(ctx, legacyParams)
	return nil
}
