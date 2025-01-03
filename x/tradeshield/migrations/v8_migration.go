package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {
	allSpotOrders := m.keeper.GetAllLegacyPendingSpotOrder(ctx)
	for _, order := range allSpotOrders {
		m.keeper.RemovePendingSpotOrder(ctx, order.OrderId)
		m.keeper.SetPendingSpotOrder(ctx, types.SpotOrder{
			OrderType:        order.OrderType,
			OrderId:          order.OrderId,
			OrderPrice:       order.OrderPrice.Rate,
			OrderAmount:      order.OrderAmount,
			OwnerAddress:     order.OwnerAddress,
			OrderTargetDenom: order.OrderTargetDenom,
			Status:           order.Status,
			Date:             order.Date,
		})
	}

	allPerpOrders := m.keeper.GetAllLegacyPendingPerpetualOrder(ctx)
	for _, order := range allPerpOrders {
		m.keeper.RemovePendingPerpetualOrder(ctx, order.OrderId)
		m.keeper.SetPendingPerpetualOrder(ctx, types.PerpetualOrder{
			OrderId:            order.OrderId,
			OwnerAddress:       order.OwnerAddress,
			PerpetualOrderType: order.PerpetualOrderType,
			Position:           order.Position,
			TriggerPrice:       order.TriggerPrice.Rate,
			Collateral:         order.Collateral,
			TradingAsset:       order.TradingAsset,
			Leverage:           order.Leverage,
			TakeProfitPrice:    order.TakeProfitPrice,
			PositionId:         order.PositionId,
			Status:             order.Status,
			StopLossPrice:      order.StopLossPrice,
			PoolId:             order.PoolId,
		})
	}

	return nil
}
