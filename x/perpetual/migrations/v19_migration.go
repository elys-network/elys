package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V19Migration(ctx sdk.Context) error {

	params := m.keeper.GetParams(ctx)
	params.PoolMaxLiabilitiesThreshold = math.LegacyMustNewDecFromStr("0.3")

	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}

	err = m.keeper.ResetStore(ctx)
	if err != nil {
		return err
	}
	return nil
}
