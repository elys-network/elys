package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (m Migrator) V15Migration(ctx sdk.Context) error {
	legacyPools := m.keeper.GetAllLegacyPools(ctx)

	for _, legacyPool := range legacyPools {
		newPool := types.Pool{
			AmmPoolId:          legacyPool.AmmPoolId,
			Health:             legacyPool.Health,
			LeveragedLpAmount:  legacyPool.LeveragedLpAmount,
			LeverageMax:        legacyPool.LeverageMax,
			MaxLeveragelpRatio: math.LegacyMustNewDecFromStr("0.6"),
		}
		m.keeper.SetPool(ctx, newPool)
	}

	return nil
}
