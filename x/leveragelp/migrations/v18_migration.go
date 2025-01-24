package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V18Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetParams(ctx)
	legacyParams.ExitFee = math.LegacyMustNewDecFromStr("0.002")

	err := m.keeper.SetParams(ctx, &legacyParams)
	if err != nil {
		return err
	}
	return nil
}
