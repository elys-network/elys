package migrations

import (
<<<<<<< HEAD
=======
	"cosmossdk.io/math"
>>>>>>> 267bed94a9ef69af6b2214edf6bf602090c98a11
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {
<<<<<<< HEAD
	m.keeper.MoveAllDebt(ctx)
	m.keeper.MoveAllInterest(ctx)
=======
	params := m.keeper.GetParams(ctx)
	params.MaxWithdrawRatio = math.LegacyMustNewDecFromStr("0.9")
	m.keeper.SetParams(ctx, params)
>>>>>>> 267bed94a9ef69af6b2214edf6bf602090c98a11
	return nil
}
