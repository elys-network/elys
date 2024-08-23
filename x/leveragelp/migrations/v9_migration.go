package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (m Migrator) V9Migration(ctx sdk.Context) error {
	positions := m.keeper.GetAllLegacyPositions(ctx)

	ctx.Logger().Info("Migrating positions from legacy to new format")

	for _, position := range positions {
		new_position := types.Position{
			Address:           position.Address,
			Collateral:        position.Collateral,
			Liabilities:       position.Liabilities,
			Leverage:          position.Leverage,
			LeveragedLpAmount: position.LeveragedLpAmount,
			PositionHealth:    position.PositionHealth,
			Id:                position.Id,
			AmmPoolId:         position.AmmPoolId,
			StopLossPrice:     position.StopLossPrice,
		}
		m.keeper.DeleteLegacyPosition(ctx, position.Address, position.Id)
		m.keeper.SetPosition(ctx, &new_position)
	}

	return nil
}
