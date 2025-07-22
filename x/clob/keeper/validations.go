package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
)

// minDec returns the minimum of two math.LegacyDec values
func minDec(a, b math.LegacyDec) math.LegacyDec {
	if a.LT(b) {
		return a
	}
	return b
}

// ValidateOrderPrice validates that an order price is within acceptable bounds
func (k Keeper) ValidateOrderPrice(ctx sdk.Context, market types.PerpetualMarket, orderType types.OrderType, price math.LegacyDec) error {
	// Price must be positive
	if !price.IsPositive() {
		return fmt.Errorf("order price must be positive, got: %s", price)
	}

	// Get current market prices for validation
	midPrice, err := k.GetMidPrice(ctx, market.Id)
	if err != nil {
		// If no mid price, check for any price reference
		bestBid := k.GetHighestBuyPrice(ctx, market.Id)
		bestAsk := k.GetLowestSellPrice(ctx, market.Id)

		if bestBid.IsZero() && bestAsk.IsZero() {
			// No price reference, accept any reasonable price
			return nil
		}

		if !bestBid.IsZero() && !bestAsk.IsZero() {
			midPrice = bestBid.Add(bestAsk).QuoInt64(2)
		} else if !bestBid.IsZero() {
			midPrice = bestBid
		} else {
			midPrice = bestAsk
		}
	}

	// Define price bounds (e.g., 50% deviation from mid price)
	maxPriceDeviation := math.LegacyNewDecWithPrec(50, 2) // 50%
	upperBound := midPrice.Mul(math.LegacyOneDec().Add(maxPriceDeviation))
	lowerBound := midPrice.Mul(math.LegacyOneDec().Sub(maxPriceDeviation))

	// Validate price is within bounds
	if price.GT(upperBound) {
		return fmt.Errorf("order price %s exceeds maximum allowed price %s (50%% above mid price %s)",
			price, upperBound, midPrice)
	}
	if price.LT(lowerBound) {
		return fmt.Errorf("order price %s is below minimum allowed price %s (50%% below mid price %s)",
			price, lowerBound, midPrice)
	}

	// Additional validation for limit orders
	if orderType == types.OrderType_ORDER_TYPE_LIMIT_BUY {
		// Buy limit orders should not be too far above current ask
		bestAsk := k.GetLowestSellPrice(ctx, market.Id)
		if !bestAsk.IsZero() && price.GT(bestAsk.Mul(math.LegacyNewDecWithPrec(120, 2))) {
			return fmt.Errorf("buy limit order price %s is too far above best ask %s", price, bestAsk)
		}
	} else if orderType == types.OrderType_ORDER_TYPE_LIMIT_SELL {
		// Sell limit orders should not be too far below current bid
		bestBid := k.GetHighestBuyPrice(ctx, market.Id)
		if !bestBid.IsZero() && price.LT(bestBid.Mul(math.LegacyNewDecWithPrec(80, 2))) {
			return fmt.Errorf("sell limit order price %s is too far below best bid %s", price, bestBid)
		}
	}

	return nil
}

// ValidateOrderQuantity validates that an order quantity is within acceptable bounds
func (k Keeper) ValidateOrderQuantity(ctx sdk.Context, market types.PerpetualMarket, quantity math.LegacyDec) error {
	// Quantity must be positive
	if !quantity.IsPositive() {
		return fmt.Errorf("order quantity must be positive, got: %s", quantity)
	}

	// Check minimum quantity tick size
	if market.MinQuantityTickSize.IsPositive() && quantity.LT(market.MinQuantityTickSize) {
		return fmt.Errorf("order quantity %s is below minimum tick size %s", quantity, market.MinQuantityTickSize)
	}

	// Check quantity is a multiple of tick size
	if market.MinQuantityTickSize.IsPositive() && !quantity.Quo(market.MinQuantityTickSize).IsInteger() {
		return fmt.Errorf("order quantity %s is not a multiple of tick size %s", quantity, market.MinQuantityTickSize)
	}

	// Check against market's total open interest limit if available
	// Note: This is a placeholder check - implement actual max open interest logic based on your requirements
	maxOpenInterest := market.InitialMarginRatio.MulInt64(1000000) // Example limit based on margin ratio

	if quantity.GT(maxOpenInterest) {
		return fmt.Errorf("order quantity %s exceeds maximum allowed %s", quantity, maxOpenInterest)
	}

	return nil
}

// ValidateSlippageProtection validates slippage protection for market orders
func (k Keeper) ValidateSlippageProtection(ctx sdk.Context, market types.PerpetualMarket,
	orderType types.OrderType, quantity math.LegacyDec, worstPrice *math.LegacyDec) error {

	if worstPrice == nil || worstPrice.IsZero() {
		// No slippage protection requested
		return nil
	}

	// Calculate expected execution price based on order book
	var expectedPrice math.LegacyDec
	var found bool

	if orderType == types.OrderType_ORDER_TYPE_MARKET_BUY {
		// For buy orders, check sell order book
		sellOrders := k.GetSellOrdersUpToQuantity(ctx, market.Id, quantity)
		if len(sellOrders) == 0 {
			return fmt.Errorf("insufficient liquidity for market buy order")
		}

		// Calculate volume-weighted average price
		totalValue := math.LegacyZeroDec()
		totalQuantity := math.LegacyZeroDec()

		for _, order := range sellOrders {
			orderQty := minDec(order.Quantity, quantity.Sub(totalQuantity))
			totalValue = totalValue.Add(order.Price.Mul(orderQty))
			totalQuantity = totalQuantity.Add(orderQty)

			if totalQuantity.GTE(quantity) {
				break
			}
		}

		if totalQuantity.LT(quantity) {
			return fmt.Errorf("insufficient liquidity: only %s available, need %s", totalQuantity, quantity)
		}

		expectedPrice = totalValue.Quo(totalQuantity)
		found = true

		// For buy orders, worst price is the maximum acceptable price
		if found && expectedPrice.GT(*worstPrice) {
			return fmt.Errorf("expected execution price %s exceeds worst acceptable price %s",
				expectedPrice, *worstPrice)
		}

	} else if orderType == types.OrderType_ORDER_TYPE_MARKET_SELL {
		// For sell orders, check buy order book
		buyOrders := k.GetBuyOrdersUpToQuantity(ctx, market.Id, quantity)
		if len(buyOrders) == 0 {
			return fmt.Errorf("insufficient liquidity for market sell order")
		}

		// Calculate volume-weighted average price
		totalValue := math.LegacyZeroDec()
		totalQuantity := math.LegacyZeroDec()

		for _, order := range buyOrders {
			orderQty := minDec(order.Quantity, quantity.Sub(totalQuantity))
			totalValue = totalValue.Add(order.Price.Mul(orderQty))
			totalQuantity = totalQuantity.Add(orderQty)

			if totalQuantity.GTE(quantity) {
				break
			}
		}

		if totalQuantity.LT(quantity) {
			return fmt.Errorf("insufficient liquidity: only %s available, need %s", totalQuantity, quantity)
		}

		expectedPrice = totalValue.Quo(totalQuantity)
		found = true

		// For sell orders, worst price is the minimum acceptable price
		if found && expectedPrice.LT(*worstPrice) {
			return fmt.Errorf("expected execution price %s is below worst acceptable price %s",
				expectedPrice, *worstPrice)
		}
	}

	return nil
}

// ValidatePositionLeverage validates that a position's leverage is within acceptable bounds
func (k Keeper) ValidatePositionLeverage(ctx sdk.Context, market types.PerpetualMarket,
	account sdk.AccAddress, additionalMargin math.Int, positionSize math.LegacyDec) error {

	// Get account's total collateral
	subAccount, err := k.GetSubAccount(ctx, account, market.Id)
	if err != nil {
		return err
	}

	totalCollateral := k.GetSubAccountBalanceOf(ctx, subAccount, market.QuoteDenom).Amount
	totalCollateral = totalCollateral.Add(additionalMargin)

	if !totalCollateral.IsPositive() {
		return fmt.Errorf("insufficient collateral")
	}

	// Calculate position value
	markPrice, err := k.GetAssetPriceFromDenom(ctx, market.BaseDenom)
	if err != nil {
		return err
	}

	positionValue := positionSize.Abs().Mul(markPrice)

	// Calculate required initial margin
	requiredMargin := positionValue.Mul(market.InitialMarginRatio)

	// Check if account has sufficient margin
	if totalCollateral.ToLegacyDec().LT(requiredMargin) {
		leverage := positionValue.Quo(totalCollateral.ToLegacyDec())
		maxLeverage := math.LegacyOneDec().Quo(market.InitialMarginRatio)
		return fmt.Errorf("insufficient margin: position leverage %s exceeds maximum allowed %s",
			leverage, maxLeverage)
	}

	return nil
}

// GetSellOrdersUpToQuantity returns sell orders up to the specified quantity
func (k Keeper) GetSellOrdersUpToQuantity(ctx sdk.Context, marketId uint64, quantity math.LegacyDec) []types.OrderBookEntry {
	sellIterator := k.GetSellOrderIterator(ctx, marketId)
	defer sellIterator.Close()

	var orders []types.OrderBookEntry
	totalQuantity := math.LegacyZeroDec()

	for ; sellIterator.Valid() && totalQuantity.LT(quantity); sellIterator.Next() {
		var order types.PerpetualOrder
		if err := k.cdc.Unmarshal(sellIterator.Value(), &order); err != nil {
			ctx.Logger().Error("failed to unmarshal sell order", "error", err)
			continue
		}

		remainingQty := order.Amount.Sub(order.Filled)
		orders = append(orders, types.OrderBookEntry{
			Price:    order.Price,
			Quantity: remainingQty,
		})
		totalQuantity = totalQuantity.Add(remainingQty)
	}

	return orders
}

// GetBuyOrdersUpToQuantity returns buy orders up to the specified quantity
func (k Keeper) GetBuyOrdersUpToQuantity(ctx sdk.Context, marketId uint64, quantity math.LegacyDec) []types.OrderBookEntry {
	buyIterator := k.GetBuyOrderIterator(ctx, marketId)
	defer buyIterator.Close()

	var orders []types.OrderBookEntry
	totalQuantity := math.LegacyZeroDec()

	// For buy orders, we want the highest prices first
	var allOrders []types.PerpetualOrder
	for ; buyIterator.Valid(); buyIterator.Next() {
		var order types.PerpetualOrder
		if err := k.cdc.Unmarshal(buyIterator.Value(), &order); err != nil {
			ctx.Logger().Error("failed to unmarshal buy order", "error", err)
			continue
		}
		allOrders = append(allOrders, order)
	}

	// Process in reverse order (highest price first)
	for i := len(allOrders) - 1; i >= 0 && totalQuantity.LT(quantity); i-- {
		order := allOrders[i]
		remainingQty := order.Amount.Sub(order.Filled)
		orders = append(orders, types.OrderBookEntry{
			Price:    order.Price,
			Quantity: remainingQty,
		})
		totalQuantity = totalQuantity.Add(remainingQty)
	}

	return orders
}
