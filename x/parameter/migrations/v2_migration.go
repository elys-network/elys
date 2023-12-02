package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	params := types.NewParams(
		sdk.NewDecWithPrec(5, 2),  // min commission 0.05
		sdk.NewDecWithPrec(66, 1), // max voting power
		sdk.NewInt(1),             // min self delegation
		"elys1mx32w9tnfxv0z5j000750h8ver7qf3xpj09w3uzvsr3hq68f4hxqte4gam", // broker address
	)
	m.keeper.SetParams(ctx, params)
	return nil
}
