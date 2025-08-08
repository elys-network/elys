package keeper

import (
	"sort"
	"sync"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

type OrderList struct {
	Orders []types.Order
}

// MemoryOrderBook maintains in-memory order book synchronized with blockchain state
type MemoryOrderBook struct {
	mu sync.RWMutex

	// Map structure: marketId -> priceTick -> orders at that price level
	// Orders within each price level are sorted by counter (FIFO)
	buyOrdersByPrice  map[uint64]map[types.PriceTick]*OrderList // marketId -> price -> orders
	sellOrdersByPrice map[uint64]map[types.PriceTick]*OrderList // marketId -> price -> orders

	// Track sorted price levels for efficient iteration
	buyPriceLevels  map[uint64][]types.PriceTick // sorted descending
	sellPriceLevels map[uint64][]types.PriceTick // sorted ascending

	// Track orders by ID for quick lookup/updates
	ordersByID map[uint64]map[uint64]*types.Order // marketId -> orderId -> order

	initialized bool
}

func NewMemoryOrderBook() *MemoryOrderBook {
	return &MemoryOrderBook{
		buyOrdersByPrice:  make(map[uint64]map[types.PriceTick]*OrderList),
		sellOrdersByPrice: make(map[uint64]map[types.PriceTick]*OrderList),
		buyPriceLevels:    make(map[uint64][]types.PriceTick),
		sellPriceLevels:   make(map[uint64][]types.PriceTick),
		ordersByID:        make(map[uint64]map[uint64]*types.Order),
		initialized:       false,
	}
}

// insertBuyPriceLevel inserts a buy price level maintaining descending order
func (mob *MemoryOrderBook) insertBuyPriceLevel(marketId uint64, priceTick types.PriceTick) {
	levels := mob.buyPriceLevels[marketId]
	index := sort.Search(len(levels), func(i int) bool {
		return levels[i] <= priceTick
	})

	// Check if price already exists
	if index < len(levels) && levels[index] == priceTick {
		return
	}

	// Insert at the correct position
	mob.buyPriceLevels[marketId] = append(levels[:index], append([]types.PriceTick{priceTick}, levels[index:]...)...)
}

// insertSellPriceLevel inserts a sell price level maintaining ascending order
func (mob *MemoryOrderBook) insertSellPriceLevel(marketId uint64, priceTick types.PriceTick) {
	levels := mob.sellPriceLevels[marketId]
	index := sort.Search(len(levels), func(i int) bool {
		return levels[i] >= priceTick
	})

	// Check if price already exists
	if index < len(levels) && levels[index] == priceTick {
		return
	}

	// Insert at the correct position
	mob.sellPriceLevels[marketId] = append(levels[:index], append([]types.PriceTick{priceTick}, levels[index:]...)...)
}

// removeBuyPriceLevel removes a buy price level
func (mob *MemoryOrderBook) removeBuyPriceLevel(marketId uint64, priceTick types.PriceTick) {
	levels := mob.buyPriceLevels[marketId]
	index := sort.Search(len(levels), func(i int) bool {
		return levels[i] <= priceTick
	})

	if index < len(levels) && levels[index] == priceTick {
		mob.buyPriceLevels[marketId] = append(levels[:index], levels[index+1:]...)
	}
}

// removeSellPriceLevel removes a sell price level
func (mob *MemoryOrderBook) removeSellPriceLevel(marketId uint64, priceTick types.PriceTick) {
	levels := mob.sellPriceLevels[marketId]
	index := sort.Search(len(levels), func(i int) bool {
		return levels[i] >= priceTick
	})

	if index < len(levels) && levels[index] == priceTick {
		mob.sellPriceLevels[marketId] = append(levels[:index], levels[index+1:]...)
	}
}

// insertOrderSorted inserts order into list maintaining FIFO order by counter
func insertOrderSorted(orders []types.Order, order types.Order) []types.Order {
	// Orders are sorted by counter (FIFO - first in first out)
	// Lower counter = earlier order = higher priority
	index := sort.Search(len(orders), func(i int) bool {
		return orders[i].GetCounter() > order.GetCounter()
	})

	orders = append(orders, types.Order{})
	copy(orders[index+1:], orders[index:])
	orders[index] = order
	return orders
}

// IsInitialized returns whether the memory orderbook has been initialized
func (mob *MemoryOrderBook) IsInitialized() bool {
	mob.mu.RLock()
	defer mob.mu.RUnlock()
	return mob.initialized
}

// InitializeFromState loads all orders from blockchain state on startup
func (mob *MemoryOrderBook) InitializeFromState(k Keeper, ctx sdk.Context) {
	mob.mu.Lock()
	defer mob.mu.Unlock()

	if mob.initialized {
		return
	}

	// Get all markets
	markets := k.GetAllPerpetualMarket(ctx)

	for _, market := range markets {
		marketId := market.Id

		// Initialize maps for this market
		mob.buyOrdersByPrice[marketId] = make(map[types.PriceTick]*OrderList)
		mob.sellOrdersByPrice[marketId] = make(map[types.PriceTick]*OrderList)
		mob.ordersByID[marketId] = make(map[uint64]*types.Order)
		mob.buyPriceLevels[marketId] = []types.PriceTick{}
		mob.sellPriceLevels[marketId] = []types.PriceTick{}

		// Load buy orders - iterator returns them in correct order (highest price first)
		buyIterator := k.GetBuyOrderIterator(ctx, marketId)
		for ; buyIterator.Valid(); buyIterator.Next() {
			var order types.Order
			k.cdc.MustUnmarshal(buyIterator.Value(), &order)
			mob.addOrderInternal(marketId, &order, true)
		}
		buyIterator.Close()

		// Load sell orders - iterator returns them in correct order (lowest price first)
		sellIterator := k.GetSellOrderIterator(ctx, marketId)
		for ; sellIterator.Valid(); sellIterator.Next() {
			var order types.Order
			k.cdc.MustUnmarshal(sellIterator.Value(), &order)
			mob.addOrderInternal(marketId, &order, false)
		}
		sellIterator.Close()
	}

	mob.initialized = true
}

// addOrderInternal adds order to memory without locking (must be called with lock held)
func (mob *MemoryOrderBook) addOrderInternal(marketId uint64, order *types.Order, isBuy bool) {
	priceTick := order.GetPriceTick()
	orderId := order.GetCounter()

	// Create a copy of the order to store
	orderCopy := *order

	// Store order by ID
	if mob.ordersByID[marketId] == nil {
		mob.ordersByID[marketId] = make(map[uint64]*types.Order)
	}
	mob.ordersByID[marketId][orderId] = &orderCopy

	// Add to price level
	var priceMap map[types.PriceTick]*OrderList

	if isBuy {
		if mob.buyOrdersByPrice[marketId] == nil {
			mob.buyOrdersByPrice[marketId] = make(map[types.PriceTick]*OrderList)
			mob.buyPriceLevels[marketId] = []types.PriceTick{}
		}
		priceMap = mob.buyOrdersByPrice[marketId]

		// Add price level if new
		if priceMap[priceTick] == nil {
			priceMap[priceTick] = &OrderList{Orders: []types.Order{}}
			// Insert price level in sorted order
			mob.insertBuyPriceLevel(marketId, priceTick)
		}
	} else {
		if mob.sellOrdersByPrice[marketId] == nil {
			mob.sellOrdersByPrice[marketId] = make(map[types.PriceTick]*OrderList)
			mob.sellPriceLevels[marketId] = []types.PriceTick{}
		}
		priceMap = mob.sellOrdersByPrice[marketId]

		// Add price level if new
		if priceMap[priceTick] == nil {
			priceMap[priceTick] = &OrderList{Orders: []types.Order{}}
			// Insert price level in sorted order
			mob.insertSellPriceLevel(marketId, priceTick)
		}
	}

	// Insert order maintaining FIFO order by counter
	priceMap[priceTick].Orders = insertOrderSorted(priceMap[priceTick].Orders, orderCopy)

	// Update the pointer in ordersByID to point to the actual order in the list
	for i := range priceMap[priceTick].Orders {
		if priceMap[priceTick].Orders[i].GetCounter() == orderId {
			mob.ordersByID[marketId][orderId] = &priceMap[priceTick].Orders[i]
			break
		}
	}
}

// AddOrder adds a new order to the memory orderbook
func (mob *MemoryOrderBook) AddOrder(marketId uint64, order *types.Order) {
	mob.mu.Lock()
	defer mob.mu.Unlock()

	mob.addOrderInternal(marketId, order, order.IsBuy())
}

// UpdateOrder updates an existing order in memory
func (mob *MemoryOrderBook) UpdateOrder(marketId uint64, order *types.Order) {
	mob.mu.Lock()
	defer mob.mu.Unlock()

	orderId := order.GetCounter()

	// Find existing order
	existingOrder, exists := mob.ordersByID[marketId][orderId]
	if !exists {
		// If order doesn't exist, add it as new
		mob.addOrderInternal(marketId, order, order.IsBuy())
		return
	}

	// Check if fully filled
	if order.Filled.GTE(order.Amount) {
		mob.removeOrderInternal(marketId, orderId, order.IsBuy())
		return
	}

	// Update the existing order's filled amount
	existingOrder.Filled = order.Filled
	existingOrder.Amount = order.Amount
}

// removeOrderInternal removes order from memory without locking
func (mob *MemoryOrderBook) removeOrderInternal(marketId uint64, orderId uint64, isBuy bool) {
	order, exists := mob.ordersByID[marketId][orderId]
	if !exists {
		return
	}

	priceTick := order.GetPriceTick()

	// Remove from ID map
	delete(mob.ordersByID[marketId], orderId)

	// Remove from price level
	var priceMap map[types.PriceTick]*OrderList

	if isBuy {
		priceMap = mob.buyOrdersByPrice[marketId]
	} else {
		priceMap = mob.sellOrdersByPrice[marketId]
	}

	if orderList, exists := priceMap[priceTick]; exists {
		// Remove order from list
		newOrders := []types.Order{}
		for _, o := range orderList.Orders {
			if o.GetCounter() != orderId {
				newOrders = append(newOrders, o)
			}
		}

		if len(newOrders) == 0 {
			// Remove empty price level
			delete(priceMap, priceTick)

			// Remove from sorted price levels
			if isBuy {
				mob.removeBuyPriceLevel(marketId, priceTick)
			} else {
				mob.removeSellPriceLevel(marketId, priceTick)
			}
		} else {
			orderList.Orders = newOrders
			// Update pointers in ordersByID for remaining orders
			for i := range orderList.Orders {
				if id := orderList.Orders[i].GetCounter(); mob.ordersByID[marketId][id] != nil {
					mob.ordersByID[marketId][id] = &orderList.Orders[i]
				}
			}
		}
	}
}

// RemoveOrder removes an order from the memory orderbook
func (mob *MemoryOrderBook) RemoveOrder(marketId uint64, orderId uint64, isBuy bool) {
	mob.mu.Lock()
	defer mob.mu.Unlock()

	mob.removeOrderInternal(marketId, orderId, isBuy)
}

// MatchOrders performs order matching for a market using price levels
func (mob *MemoryOrderBook) MatchOrders(marketId uint64, maxMatches int) []MatchedOrder {
	mob.mu.RLock()
	defer mob.mu.RUnlock()

	matched := []MatchedOrder{}

	buyLevels := mob.buyPriceLevels[marketId]
	sellLevels := mob.sellPriceLevels[marketId]

	if len(buyLevels) == 0 || len(sellLevels) == 0 {
		return matched
	}

	buyIndex := 0
	sellIndex := 0

	for buyIndex < len(buyLevels) && sellIndex < len(sellLevels) && len(matched) < maxMatches {
		buyPrice := buyLevels[buyIndex]
		sellPrice := sellLevels[sellIndex]

		// Check if prices cross
		if buyPrice < sellPrice {
			break // No more matches possible
		}

		// Get orders at these price levels
		buyOrderList := mob.buyOrdersByPrice[marketId][buyPrice]
		sellOrderList := mob.sellOrdersByPrice[marketId][sellPrice]

		if buyOrderList == nil || sellOrderList == nil {
			break
		}

		// Match orders within these price levels (FIFO - orders are already sorted by counter)
		buyOrderIndex := 0
		sellOrderIndex := 0

		for buyOrderIndex < len(buyOrderList.Orders) &&
			sellOrderIndex < len(sellOrderList.Orders) &&
			len(matched) < maxMatches {

			buyOrder := &buyOrderList.Orders[buyOrderIndex]
			sellOrder := &sellOrderList.Orders[sellOrderIndex]

			// Skip if either order is already fully filled
			buyRemaining := buyOrder.Amount.Sub(buyOrder.Filled)
			sellRemaining := sellOrder.Amount.Sub(sellOrder.Filled)

			if !buyRemaining.IsPositive() {
				buyOrderIndex++
				continue
			}
			if !sellRemaining.IsPositive() {
				sellOrderIndex++
				continue
			}

			// Determine trade price based on order priority (counter = time priority)
			tradePrice := sellOrder.GetPrice()
			if sellOrder.GetCounter() > buyOrder.GetCounter() {
				tradePrice = buyOrder.GetPrice()
			}

			// Calculate trade quantity
			tradeQuantity := math.LegacyMinDec(buyRemaining, sellRemaining)

			if tradeQuantity.IsPositive() {
				matched = append(matched, MatchedOrder{
					BuyOrder:  *buyOrder,
					SellOrder: *sellOrder,
					Price:     tradePrice,
					Quantity:  tradeQuantity,
					MarketId:  marketId,
				})
			}

			// Move to next order based on which one would be fully filled
			if buyRemaining.LTE(sellRemaining) {
				buyOrderIndex++
			}
			if sellRemaining.LTE(buyRemaining) {
				sellOrderIndex++
			}
		}

		// Move to next price level if all orders at current level are processed
		if buyOrderIndex >= len(buyOrderList.Orders) {
			buyIndex++
		}
		if sellOrderIndex >= len(sellOrderList.Orders) {
			sellIndex++
		}
	}

	return matched
}

type MatchedOrder struct {
	BuyOrder  types.Order
	SellOrder types.Order
	Price     math.LegacyDec
	Quantity  math.LegacyDec
	MarketId  uint64
}

// GetOperationsToPropose returns matched orders for vote extension
func (k Keeper) GetOperationsToPropose(ctx sdk.Context, marketIds []uint64, maxMatchesPerMarket int) []types.MatchedOrderExecution {
	// Ensure memory orderbook is initialized
	if !k.memoryOrderBook.initialized {
		k.memoryOrderBook.InitializeFromState(k, ctx)
	}

	allMatched := []types.MatchedOrderExecution{}

	for _, marketId := range marketIds {
		// Match orders for this market
		matched := k.memoryOrderBook.MatchOrders(marketId, maxMatchesPerMarket)

		// Convert to MatchedOrderExecution format for vote extension
		for _, m := range matched {
			allMatched = append(allMatched, types.MatchedOrderExecution{
				BuyOrderCounter:    m.BuyOrder.GetCounter(),
				SellOrderCounter:   m.SellOrder.GetCounter(),
				MarketId:           m.MarketId,
				Price:              m.Price,
				Quantity:           m.Quantity,
				Buyer:              m.BuyOrder.Owner,
				Seller:             m.SellOrder.Owner,
				BuyerSubAccountId:  m.BuyOrder.SubAccountId,
				SellerSubAccountId: m.SellOrder.SubAccountId,
			})
		}
	}

	return allMatched
}
