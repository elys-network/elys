package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V20Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.ExitBuffer = math.LegacyMustNewDecFromStr("0.1")
	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}
	return nil
}
