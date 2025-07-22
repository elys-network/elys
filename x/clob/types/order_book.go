package types

import (
	"cosmossdk.io/math"
)

// OrderBookEntry represents a single level in the order book
type OrderBookEntry struct {
	Price    math.LegacyDec `json:"price"`
	Quantity math.LegacyDec `json:"quantity"`
}

// OrderBookSnapshot represents a snapshot of the order book at a specific time
type OrderBookSnapshot struct {
	MarketId   uint64           `json:"market_id"`
	BuyOrders  []OrderBookEntry `json:"buy_orders"`
	SellOrders []OrderBookEntry `json:"sell_orders"`
	Timestamp  uint64           `json:"timestamp"`
}

// OrderBookUpdate represents an update to the order book
type OrderBookUpdate struct {
	MarketId  uint64         `json:"market_id"`
	Side      OrderSide      `json:"side"`
	Operation UpdateOp       `json:"operation"`
	Price     math.LegacyDec `json:"price"`
	Quantity  math.LegacyDec `json:"quantity"`
}

// OrderSide represents the side of an order
type OrderSide int32

const (
	OrderSideBuy OrderSide = iota
	OrderSideSell
)

// UpdateOp represents the type of update operation
type UpdateOp int32

const (
	UpdateOpAdd UpdateOp = iota
	UpdateOpRemove
	UpdateOpModify
)

// String returns the string representation of OrderSide
func (s OrderSide) String() string {
	switch s {
	case OrderSideBuy:
		return "buy"
	case OrderSideSell:
		return "sell"
	default:
		return "unknown"
	}
}

// String returns the string representation of UpdateOp
func (o UpdateOp) String() string {
	switch o {
	case UpdateOpAdd:
		return "add"
	case UpdateOpRemove:
		return "remove"
	case UpdateOpModify:
		return "modify"
	default:
		return "unknown"
	}
}
