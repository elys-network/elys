package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.MaxWithdrawRatio = math.LegacyMustNewDecFromStr("0.9")
	m.keeper.SetParams(ctx, params)
	return nil
}
