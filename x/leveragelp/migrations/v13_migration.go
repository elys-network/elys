package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (m Migrator) V13Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllLegacyPools(ctx)

	for _, pool := range pools {
		new_pool := types.Pool{
			AmmPoolId: pool.AmmPoolId,
			Health: pool.Health,
			Enabled: pool.Enabled,
			Closed: pool.Closed,
			LeveragedLpAmount: pool.LeveragedLpAmount,
			LeverageMax: pool.LeverageMax,
			MaxLeveragelpPercent: sdk.NewDec(60),
		}
		m.keeper.SetPool(ctx, new_pool)
		m.keeper.DeleteLegacyPool(ctx, pool.AmmPoolId)
	}

	return nil
}