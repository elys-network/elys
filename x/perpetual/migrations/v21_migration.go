package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V21Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.TakerFee = math.LegacyMustNewDecFromStr("0.0005")
	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}
	return nil
}
