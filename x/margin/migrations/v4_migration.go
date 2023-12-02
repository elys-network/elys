package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
	// reset params
	params := types.NewParams()
	m.keeper.SetParams(ctx, &params)

	// reset mtps
	for _, mtp := range m.keeper.GetAllMTPs(ctx) {
		m.keeper.DestroyMTP(ctx, mtp.Address, mtp.Id)
	}

	// reset pools
	for _, pool := range m.keeper.GetAllPools(ctx) {
		m.keeper.RemovePool(ctx, pool.AmmPoolId)
	}

	return nil
}
