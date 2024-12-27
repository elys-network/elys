package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {

	// Testnet
	if ctx.ChainID() == "elysicstestnet-1" {
		allOrders := m.keeper.GetAllPendingSpotOrder(ctx)
		for _, order := range allOrders {
			if order.OrderType == types.SpotOrderType_LIMITBUY {
				m.keeper.RemovePendingSpotOrder(ctx, order.OrderId)
			}
		}
	} else { // Mainnet
		ordersIds := []uint64{164, 177, 180}
		for _, orderId := range ordersIds {
			m.keeper.RemovePendingSpotOrder(ctx, orderId)
		}
	}

	return nil
}
