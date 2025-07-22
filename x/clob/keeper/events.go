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
			sdk.NewAttribute("buyer", trade.BuyerSubAccount.Owner),
			sdk.NewAttribute("seller", trade.SellerSubAccount.Owner),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", trade.MarketId)),
			sdk.NewAttribute(types.AttributePrice, trade.Price.String()),
			sdk.NewAttribute(types.AttributeQuantity, trade.Quantity.String()),
			sdk.NewAttribute("is_buyer_liquidation", fmt.Sprintf("%t", trade.IsBuyerLiquidation)),
			sdk.NewAttribute("is_seller_liquidation", fmt.Sprintf("%t", trade.IsSellerLiquidation)),
			sdk.NewAttribute("is_buyer_taker", fmt.Sprintf("%t", trade.IsBuyerTaker)),
		),
	)
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
			sdk.NewAttribute("funding_rate", fundingRate.String()),
		),
	)
}
