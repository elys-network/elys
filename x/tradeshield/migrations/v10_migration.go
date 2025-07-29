package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
)

func (m Migrator) V10Migration(ctx sdk.Context) error {
	if ctx.ChainID() == "elysicstestnet-1" {
		m.keeper.DeleteAllPendingPerpetualOrder(ctx)
	} else {
		legacyOrders := m.keeper.GetAllLegacyPendingPerpetualOrder(ctx)
		for _, order := range legacyOrders {
			newOrder := types.PerpetualOrder{
				OrderId:            order.OrderId,
				OwnerAddress:       order.OwnerAddress,
				PerpetualOrderType: order.PerpetualOrderType,
				Position:           order.Position,
				TriggerPrice:       order.TriggerPrice,
				Collateral:         order.Collateral,
				Leverage:           order.Leverage,
				TakeProfitPrice:    order.TakeProfitPrice,
				PositionId:         order.PositionId,
				Status:             order.Status,
				StopLossPrice:      order.StopLossPrice,
				PoolId:             order.PoolId,
			}
			m.keeper.SetPendingPerpetualOrder(ctx, newOrder)
		}
	}
	return nil
}
