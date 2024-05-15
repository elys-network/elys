package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	// reset params
	params := types.NewParams(
		sdk.NewDecWithPrec(5, 2),                      // min commission 0.05
		sdk.NewDecWithPrec(66, 1),                     // max voting power
		sdk.NewInt(1),                                 // min self delegation
		"elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l", // broker address
		6307200, // total blocks per year
	)
	m.keeper.SetParams(ctx, params)

	return nil
}
