package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
)

// EmitPositionEvent emits an event when a position is opened, closed, or modified
func (k Keeper) EmitPositionEvent(ctx sdk.Context, eventType string, perpetual types.Perpetual,
	tradePrice, tradeQuantity math.LegacyDec, pnl math.Int) {

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeOwner, perpetual.Owner),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", perpetual.MarketId)),
			sdk.NewAttribute(types.AttributePositionId, fmt.Sprintf("%d", perpetual.Id)),
			sdk.NewAttribute(types.AttributeQuantity, perpetual.Quantity.String()),
			sdk.NewAttribute(types.AttributeTradePrice, tradePrice.String()),
			sdk.NewAttribute(types.AttributeTradeQuantity, tradeQuantity.String()),
			sdk.NewAttribute(types.AttributePnL, pnl.String()),
		),
	)
}

// EmitOrderExecutedEvent emits an event when an order is executed
func (k Keeper) EmitOrderExecutedEvent(ctx sdk.Context, trade types.Trade) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventOrderExecuted,
			sdk.NewAttribute(types.AttributeBuyer, trade.BuyerSubAccount.Owner),
			sdk.NewAttribute(types.AttributeSeller, trade.SellerSubAccount.Owner),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", trade.MarketId)),
			sdk.NewAttribute(types.AttributePrice, trade.Price.String()),
			sdk.NewAttribute(types.AttributeQuantity, trade.Quantity.String()),
			sdk.NewAttribute(types.AttributeIsLiquidation, fmt.Sprintf("%t", trade.IsBuyerLiquidation || trade.IsSellerLiquidation)),
			sdk.NewAttribute(types.AttributeIsTaker, fmt.Sprintf("%t", trade.IsBuyerTaker)),
		),
	)

	// Also emit a trade event
	k.EmitTradeEvent(ctx, trade)
}

// EmitFundingPaymentEvent emits an event when funding payment is made
func (k Keeper) EmitFundingPaymentEvent(ctx sdk.Context, perpetual types.Perpetual,
	fundingAmount math.Int, fundingRate math.LegacyDec) {

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventFundingPayment,
			sdk.NewAttribute(types.AttributeOwner, perpetual.Owner),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", perpetual.MarketId)),
			sdk.NewAttribute(types.AttributePositionId, fmt.Sprintf("%d", perpetual.Id)),
			sdk.NewAttribute(types.AttributeFundingAmount, fundingAmount.String()),
			sdk.NewAttribute(types.AttributeFundingRate, fundingRate.String()),
		),
	)
}

// EmitTradeEvent emits a detailed trade event
func (k Keeper) EmitTradeEvent(ctx sdk.Context, trade types.Trade) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTrade,
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", trade.MarketId)),
			sdk.NewAttribute(types.AttributeBuyer, trade.BuyerSubAccount.Owner),
			sdk.NewAttribute(types.AttributeSeller, trade.SellerSubAccount.Owner),
			sdk.NewAttribute(types.AttributePrice, trade.Price.String()),
			sdk.NewAttribute(types.AttributeQuantity, trade.Quantity.String()),
			sdk.NewAttribute(types.AttributeTimestamp, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)
}

// EmitCancelOrderEvent emits an event when an order is cancelled
func (k Keeper) EmitCancelOrderEvent(ctx sdk.Context, owner string, marketId uint64, orderId uint64) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventCancelOrder,
			sdk.NewAttribute(types.AttributeOwner, owner),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", marketId)),
			sdk.NewAttribute(types.AttributeOrderId, fmt.Sprintf("%d", orderId)),
			sdk.NewAttribute(types.AttributeTimestamp, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)
}

// EmitCancelAllOrdersEvent emits an event when all orders are cancelled
func (k Keeper) EmitCancelAllOrdersEvent(ctx sdk.Context, owner string, marketId uint64, count uint32) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventCancelAllOrders,
			sdk.NewAttribute(types.AttributeOwner, owner),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", marketId)),
			sdk.NewAttribute(types.AttributeQuantity, fmt.Sprintf("%d", count)),
			sdk.NewAttribute(types.AttributeTimestamp, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)
}

// EmitFundingRateUpdateEvent emits an event when funding rate is updated
func (k Keeper) EmitFundingRateUpdateEvent(ctx sdk.Context, marketId uint64, previousRate, newRate math.LegacyDec) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventFundingRateUpdate,
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", marketId)),
			sdk.NewAttribute(types.AttributePreviousFundingRate, previousRate.String()),
			sdk.NewAttribute(types.AttributeFundingRate, newRate.String()),
			sdk.NewAttribute(types.AttributeBlockHeight, fmt.Sprintf("%d", ctx.BlockHeight())),
			sdk.NewAttribute(types.AttributeTimestamp, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)
}

// EmitLiquidationTriggeredEvent emits an event when liquidation is triggered
func (k Keeper) EmitLiquidationTriggeredEvent(ctx sdk.Context, perpetual types.Perpetual, liquidator string, liquidationPrice math.LegacyDec) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventLiquidationTriggered,
			sdk.NewAttribute(types.AttributeLiquidator, liquidator),
			sdk.NewAttribute(types.AttributeOwner, perpetual.Owner),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", perpetual.MarketId)),
			sdk.NewAttribute(types.AttributePositionId, fmt.Sprintf("%d", perpetual.Id)),
			sdk.NewAttribute(types.AttributeLiquidationPrice, liquidationPrice.String()),
			sdk.NewAttribute(types.AttributeTimestamp, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)
}

// EmitAutoDeleveragingEvent emits an event when auto-deleveraging occurs
func (k Keeper) EmitAutoDeleveragingEvent(ctx sdk.Context, perpetual types.Perpetual, counterparty string, quantity math.LegacyDec, price math.LegacyDec) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventAutoDeleveraging,
			sdk.NewAttribute(types.AttributeOwner, perpetual.Owner),
			sdk.NewAttribute(types.AttributeADLCounterparty, counterparty),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", perpetual.MarketId)),
			sdk.NewAttribute(types.AttributePositionId, fmt.Sprintf("%d", perpetual.Id)),
			sdk.NewAttribute(types.AttributeADLQuantity, quantity.String()),
			sdk.NewAttribute(types.AttributePrice, price.String()),
			sdk.NewAttribute(types.AttributeTimestamp, fmt.Sprintf("%d", ctx.BlockTime().Unix())),
		),
	)
}
