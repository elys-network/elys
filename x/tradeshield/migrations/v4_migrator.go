package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
	// Remove all legacy pending perpetual orders
	orders := m.keeper.GetAllLegacyPendingPerpetualOrder(ctx)
	for _, order := range orders {
		m.keeper.RemovePendingPerpetualOrder(ctx, order.OrderId)
	}
	m.keeper.DeleteAllPendingPerpetualOrder(ctx)
	m.keeper.DeleteAllPendingSpotOrder(ctx)

	return nil
}
