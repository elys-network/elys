package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V12Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetParams(ctx)
	legacyParams.MinSlippage = math.LegacyMustNewDecFromStr("0.001")

	m.keeper.SetParams(ctx, legacyParams)
	return nil
}
