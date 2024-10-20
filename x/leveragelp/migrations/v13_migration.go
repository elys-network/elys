package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (m Migrator) V13Migration(ctx sdk.Context) error {
	legacyPools := m.keeper.GetAllLegacyPools(ctx)

	for _, legacyPool := range legacyPools {
		newPool := types.Pool{
			AmmPoolId:         legacyPool.AmmPoolId,
			Health:            legacyPool.Health,
			LeveragedLpAmount: legacyPool.LeveragedLpAmount,
			LeverageMax:       legacyPool.LeverageMax,
		}
		m.keeper.SetPool(ctx, newPool)
	}

	return nil
}
