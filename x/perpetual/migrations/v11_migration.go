package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	m.keeper.NukeDB(ctx)
	params := types.DefaultParams()
	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}

	// using previous values
	params.IncrementalBorrowInterestPaymentFundAddress = "elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l"
	params.ForceCloseFundAddress = "elys1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrec2l"
	return nil
}
