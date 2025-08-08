package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) ExecuteMatchedOrders(goCtx context.Context, msg *types.MsgExecuteMatchedOrders) (*types.MsgExecuteMatchedOrdersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	totalVolume := math.LegacyZeroDec()
	ordersExecuted := uint64(0)

	// Process each matched order
	for _, matchedOrder := range msg.MatchedOrders {
		// Get the market to validate it exists
		market, err := k.GetPerpetualMarket(ctx, matchedOrder.MarketId)
		if err != nil {
			return nil, err
		}

		// Get buyer and seller addresses
		buyer, err := sdk.AccAddressFromBech32(matchedOrder.Buyer)
		if err != nil {
			return nil, err
		}

		seller, err := sdk.AccAddressFromBech32(matchedOrder.Seller)
		if err != nil {
			return nil, err
		}

		// Get subaccounts
		buyerSubAccount, err := k.GetSubAccount(ctx, buyer, matchedOrder.BuyerSubAccountId)
		if err != nil {
			return nil, err
		}

		sellerSubAccount, err := k.GetSubAccount(ctx, seller, matchedOrder.SellerSubAccountId)
		if err != nil {
			return nil, err
		}

		// Determine who is the taker based on order IDs
		isBuyerTaker := matchedOrder.BuyOrderCounter > matchedOrder.SellOrderCounter

		// Execute the trade using the existing Exchange function
		err = k.Exchange(ctx, types.Trade{
			BuyerSubAccount:     buyerSubAccount,
			SellerSubAccount:    sellerSubAccount,
			MarketId:            market.Id,
			Price:               matchedOrder.Price,
			Quantity:            matchedOrder.Quantity,
			IsBuyerLiquidation:  false,
			IsSellerLiquidation: false,
			IsBuyerTaker:        isBuyerTaker,
		})
		if err != nil {
			// Log error but continue processing other orders
			ctx.Logger().Error("Failed to execute matched order",
				"market_id", matchedOrder.MarketId,
				"buy_order_id", matchedOrder.BuyOrderCounter,
				"sell_order_id", matchedOrder.SellOrderCounter,
				"error", err)
			continue
		}

		// Update metrics
		ordersExecuted++
		totalVolume = totalVolume.Add(matchedOrder.Quantity)

		// Emit trade event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTrade,
				sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", market.Id)),
				sdk.NewAttribute(types.AttributeBuyer, matchedOrder.Buyer),
				sdk.NewAttribute(types.AttributeSeller, matchedOrder.Seller),
				sdk.NewAttribute(types.AttributeTradePrice, matchedOrder.Price.String()),
				sdk.NewAttribute(types.AttributeTradeQuantity, matchedOrder.Quantity.String()),
				sdk.NewAttribute(types.AttributeIsTaker, fmt.Sprintf("%t", isBuyerTaker)),
				sdk.NewAttribute("buy_order_counter", fmt.Sprintf("%d", matchedOrder.BuyOrderCounter)),
				sdk.NewAttribute("sell_order_counter", fmt.Sprintf("%d", matchedOrder.SellOrderCounter)),
				sdk.NewAttribute("source", "vote_extension"),
			),
		)

		// Note: Order removal/update is handled internally by the Exchange function
		// through the ExecuteLimitBuyOrder flow which updates filled amounts
		// and removes fully filled orders automatically
	}

	// Emit summary event
	if ordersExecuted > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"matched_orders_executed",
				sdk.NewAttribute("orders_executed", fmt.Sprintf("%d", ordersExecuted)),
				sdk.NewAttribute("total_volume", totalVolume.String()),
			),
		)
	}

	return &types.MsgExecuteMatchedOrdersResponse{
		OrdersExecuted: ordersExecuted,
		TotalVolume:    totalVolume,
	}, nil
}
