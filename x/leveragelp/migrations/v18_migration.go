package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V18Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetParams(ctx)
	legacyParams.ExitBuffer = math.LegacyMustNewDecFromStr("0.05")
	err := m.keeper.SetParams(ctx, &legacyParams)
	if err != nil {
		return err
	}
	m.keeper.V18MigratonPoolLiabilities(ctx)
	return nil
}
