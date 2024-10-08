package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	// Remove all legacy pending perpetual orders
	orders := m.keeper.GetAllLegacyPendingPerpetualOrder(ctx)
	for _, order := range orders {
		m.keeper.RemovePendingPerpetualOrder(ctx, order.OrderId)
	}

	return nil
}
