package migrations

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V22Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetParams(ctx)
	legacyParams.LiabilitiesFactor = sdkmath.LegacyMustNewDecFromStr("1.0")
	return m.keeper.SetParams(ctx, &legacyParams)
}
