package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"cosmossdk.io/math"
)

// Order types
type OrderType string

const (
	OrderTypeLimitBuy  OrderType = "LIMIT_BUY"
	OrderTypeLimitSell OrderType = "LIMIT_SELL"
	OrderTypeMarketBuy OrderType = "MARKET_BUY"
	OrderTypeStopLoss  OrderType = "STOP_LOSS"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusExecuted  OrderStatus = "EXECUTED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusClosed    OrderStatus = "CLOSED"
)

type PositionType string

const (
	PositionTypeLong  PositionType = "LONG"
	PositionTypeShort PositionType = "SHORT"
)

type PerpetualOrderType string

const (
	PerpetualOrderTypeLimitOpen  PerpetualOrderType = "LIMIT_OPEN"
	PerpetualOrderTypeLimitClose PerpetualOrderType = "LIMIT_CLOSE"
	PerpetualOrderTypeMarket     PerpetualOrderType = "MARKET"
	PerpetualOrderTypeStopLoss   PerpetualOrderType = "STOP_LOSS"
	PerpetualOrderTypeTakeProfit PerpetualOrderType = "TAKE_PROFIT"
)

// JSONB type for storing complex data
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Spot Order model
type SpotOrder struct {
	ID               int64       `db:"id"`
	OrderID          uint64      `db:"order_id"`
	OrderType        OrderType   `db:"order_type"`
	OwnerAddress     string      `db:"owner_address"`
	OrderTargetDenom string      `db:"order_target_denom"`
	OrderPrice       JSONB       `db:"order_price"`
	OrderAmount      math.Int    `db:"order_amount"`
	Status           OrderStatus `db:"status"`
	CreatedAt        time.Time   `db:"created_at"`
	ExecutedAt       *time.Time  `db:"executed_at"`
	ClosedAt         *time.Time  `db:"closed_at"`
	BlockHeight      int64       `db:"block_height"`
	TxHash           string      `db:"tx_hash"`
}

// Perpetual Order model
type PerpetualOrder struct {
	ID                 int64              `db:"id"`
	OrderID            uint64             `db:"order_id"`
	OwnerAddress       string             `db:"owner_address"`
	PerpetualOrderType PerpetualOrderType `db:"perpetual_order_type"`
	Position           PositionType       `db:"position"`
	TriggerPrice       JSONB              `db:"trigger_price"`
	Collateral         math.Int           `db:"collateral"`
	Leverage           math.LegacyDec     `db:"leverage"`
	TakeProfitPrice    math.LegacyDec     `db:"take_profit_price"`
	StopLossPrice      math.LegacyDec     `db:"stop_loss_price"`
	PositionID         *uint64            `db:"position_id"`
	Status             OrderStatus        `db:"status"`
	CreatedAt          time.Time          `db:"created_at"`
	ExecutedAt         *time.Time         `db:"executed_at"`
	CancelledAt        *time.Time         `db:"cancelled_at"`
	BlockHeight        int64              `db:"block_height"`
	TxHash             string             `db:"tx_hash"`
}

// Perpetual Position (MTP) model
type PerpetualPosition struct {
	ID              int64           `db:"id"`
	MtpID           uint64          `db:"mtp_id"`
	OwnerAddress    string          `db:"owner_address"`
	AmmPoolID       uint64          `db:"amm_pool_id"`
	Position        PositionType    `db:"position"`
	CollateralAsset string          `db:"collateral_asset"`
	Collateral      math.Int        `db:"collateral"`
	Liabilities     math.Int        `db:"liabilities"`
	Custody         math.Int        `db:"custody"`
	MtpHealth       math.LegacyDec  `db:"mtp_health"`
	OpenPrice       math.LegacyDec  `db:"open_price"`
	ClosingPrice    *math.LegacyDec `db:"closing_price"`
	NetPnL          *math.LegacyDec `db:"net_pnl"`
	StopLossPrice   math.LegacyDec  `db:"stop_loss_price"`
	TakeProfitPrice math.LegacyDec  `db:"take_profit_price"`
	OpenedAt        time.Time       `db:"opened_at"`
	ClosedAt        *time.Time      `db:"closed_at"`
	ClosedBy        *string         `db:"closed_by"`
	CloseTrigger    *string         `db:"close_trigger"`
	BlockHeight     int64           `db:"block_height"`
	TxHash          string          `db:"tx_hash"`
}

// Trade model
type Trade struct {
	ID           int64          `db:"id"`
	TradeType    string         `db:"trade_type"`
	ReferenceID  uint64         `db:"reference_id"`
	OwnerAddress string         `db:"owner_address"`
	Asset        string         `db:"asset"`
	Amount       math.Int       `db:"amount"`
	Price        math.LegacyDec `db:"price"`
	Fees         JSONB          `db:"fees"`
	ExecutedAt   time.Time      `db:"executed_at"`
	BlockHeight  int64          `db:"block_height"`
	TxHash       string         `db:"tx_hash"`
	EventType    string         `db:"event_type"`
}

// Order Book Entry
type OrderBookEntry struct {
	Price  math.LegacyDec `json:"price"`
	Amount math.Int       `json:"amount"`
}

// Order Book Snapshot
type OrderBookSnapshot struct {
	ID           int64            `db:"id"`
	AssetPair    string           `db:"asset_pair"`
	Bids         []OrderBookEntry `db:"bids"`
	Asks         []OrderBookEntry `db:"asks"`
	SnapshotTime time.Time        `db:"snapshot_time"`
	BlockHeight  int64            `db:"block_height"`
}

// WebSocket Subscription
type WebSocketSubscription struct {
	ID               string    `db:"id"`
	ClientID         string    `db:"client_id"`
	SubscriptionType string    `db:"subscription_type"`
	Filters          JSONB     `db:"filters"`
	CreatedAt        time.Time `db:"created_at"`
	LastPing         time.Time `db:"last_ping"`
}

// Indexer State
type IndexerState struct {
	ID                  int64     `db:"id"`
	LastProcessedHeight int64     `db:"last_processed_height"`
	LastProcessedTime   time.Time `db:"last_processed_time"`
	UpdatedAt           time.Time `db:"updated_at"`
}

// WebSocket Message Types
type WSMessageType string

const (
	WSMessageTypeSubscribe   WSMessageType = "subscribe"
	WSMessageTypeUnsubscribe WSMessageType = "unsubscribe"
	WSMessageTypePing        WSMessageType = "ping"
	WSMessageTypePong        WSMessageType = "pong"
	WSMessageTypeUpdate      WSMessageType = "update"
	WSMessageTypeError       WSMessageType = "error"
)

// WebSocket Subscription Types
type WSSubscriptionType string

const (
	WSSubscriptionOrderBook  WSSubscriptionType = "order_book"
	WSSubscriptionTrades     WSSubscriptionType = "trades"
	WSSubscriptionOrders     WSSubscriptionType = "orders"
	WSSubscriptionPositions  WSSubscriptionType = "positions"
	WSSubscriptionMarketData WSSubscriptionType = "market_data"
)

// WebSocket Messages
type WSMessage struct {
	Type    WSMessageType   `json:"type"`
	Channel string          `json:"channel,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

type WSSubscribeRequest struct {
	Type    WSSubscriptionType     `json:"type"`
	Filters map[string]interface{} `json:"filters,omitempty"`
}

type WSOrderBookUpdate struct {
	AssetPair string           `json:"asset_pair"`
	Bids      []OrderBookEntry `json:"bids"`
	Asks      []OrderBookEntry `json:"asks"`
	Timestamp time.Time        `json:"timestamp"`
}

type WSTradeUpdate struct {
	Trade     *Trade    `json:"trade"`
	Timestamp time.Time `json:"timestamp"`
}

type WSOrderUpdate struct {
	OrderType string      `json:"order_type"` // "spot" or "perpetual"
	Order     interface{} `json:"order"`      // SpotOrder or PerpetualOrder
	Action    string      `json:"action"`     // "created", "updated", "executed", "cancelled"
	Timestamp time.Time   `json:"timestamp"`
}

type WSPositionUpdate struct {
	Position  *PerpetualPosition `json:"position"`
	Action    string             `json:"action"` // "opened", "updated", "closed"
	Timestamp time.Time          `json:"timestamp"`
}
