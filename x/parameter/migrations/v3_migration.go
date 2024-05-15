package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	// reset params
	params := types.NewParams(
		sdk.NewDecWithPrec(5, 2),  // min commission 0.05
		sdk.NewDecWithPrec(66, 2), // max voting power 0.66
		sdk.NewInt(1),             // min self delegation
		"elys1m3hduhk4uzxn8mxuvpz02ysndxfwgy5mq60h4c34qqn67xud584qeee3m4", // broker address
		6307200, // total blocks per year
	)
	m.keeper.SetParams(ctx, params)

	return nil
}
