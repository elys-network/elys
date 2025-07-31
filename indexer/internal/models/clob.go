package models

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

// CLOB Order Types
type CLOBOrderType string

const (
	CLOBOrderTypeLimitBuy   CLOBOrderType = "LIMIT_BUY"
	CLOBOrderTypeLimitSell  CLOBOrderType = "LIMIT_SELL"
	CLOBOrderTypeMarketBuy  CLOBOrderType = "MARKET_BUY"
	CLOBOrderTypeMarketSell CLOBOrderType = "MARKET_SELL"
)

// CLOB Order Status
type CLOBOrderStatus string

const (
	CLOBOrderStatusPending         CLOBOrderStatus = "PENDING"
	CLOBOrderStatusPartiallyFilled CLOBOrderStatus = "PARTIALLY_FILLED"
	CLOBOrderStatusFilled          CLOBOrderStatus = "FILLED"
	CLOBOrderStatusCancelled       CLOBOrderStatus = "CANCELLED"
	CLOBOrderStatusExpired         CLOBOrderStatus = "EXPIRED"
)

// CLOB Position Side
type CLOBPositionSide string

const (
	CLOBPositionSideLong  CLOBPositionSide = "LONG"
	CLOBPositionSideShort CLOBPositionSide = "SHORT"
)

// CLOBMarket represents a CLOB perpetual market
type CLOBMarket struct {
	ID                       int64           `db:"id" json:"id"`
	MarketID                 uint64          `db:"market_id" json:"market_id"`
	Ticker                   string          `db:"ticker" json:"ticker"`
	BaseAsset                string          `db:"base_asset" json:"base_asset"`
	QuoteAsset               string          `db:"quote_asset" json:"quote_asset"`
	TickSize                 decimal.Decimal `db:"tick_size" json:"tick_size"`
	LotSize                  decimal.Decimal `db:"lot_size" json:"lot_size"`
	MinOrderSize             decimal.Decimal `db:"min_order_size" json:"min_order_size"`
	MaxOrderSize             decimal.Decimal `db:"max_order_size" json:"max_order_size"`
	MaxLeverage              decimal.Decimal `db:"max_leverage" json:"max_leverage"`
	InitialMarginFraction    decimal.Decimal `db:"initial_margin_fraction" json:"initial_margin_fraction"`
	MaintenanceMarginFraction decimal.Decimal `db:"maintenance_margin_fraction" json:"maintenance_margin_fraction"`
	FundingInterval          int64           `db:"funding_interval" json:"funding_interval"`
	NextFundingTime          time.Time       `db:"next_funding_time" json:"next_funding_time"`
	IsActive                 bool            `db:"is_active" json:"is_active"`
	CreatedAt                time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt                time.Time       `db:"updated_at" json:"updated_at"`
	BlockHeight              int64           `db:"block_height" json:"block_height"`
}

// CLOBOrder represents a CLOB order
type CLOBOrder struct {
	ID              int64             `db:"id" json:"id"`
	OrderID         uint64            `db:"order_id" json:"order_id"`
	MarketID        uint64            `db:"market_id" json:"market_id"`
	Counter         uint64            `db:"counter" json:"counter"`
	Owner           string            `db:"owner" json:"owner"`
	SubAccountID    uint64            `db:"sub_account_id" json:"sub_account_id"`
	OrderType       CLOBOrderType     `db:"order_type" json:"order_type"`
	Price           decimal.Decimal   `db:"price" json:"price"`
	Amount          decimal.Decimal   `db:"amount" json:"amount"`
	FilledAmount    decimal.Decimal   `db:"filled_amount" json:"filled_amount"`
	RemainingAmount decimal.Decimal   `db:"remaining_amount" json:"remaining_amount"`
	Status          CLOBOrderStatus   `db:"status" json:"status"`
	TimeInForce     string            `db:"time_in_force" json:"time_in_force"`
	PostOnly        bool              `db:"post_only" json:"post_only"`
	ReduceOnly      bool              `db:"reduce_only" json:"reduce_only"`
	CreatedAt       time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time         `db:"updated_at" json:"updated_at"`
	ExecutedAt      *time.Time        `db:"executed_at" json:"executed_at,omitempty"`
	CancelledAt     *time.Time        `db:"cancelled_at" json:"cancelled_at,omitempty"`
	ExpiresAt       *time.Time        `db:"expires_at" json:"expires_at,omitempty"`
	BlockHeight     int64             `db:"block_height" json:"block_height"`
	TxHash          string            `db:"tx_hash" json:"tx_hash"`
}

// CLOBTrade represents a CLOB trade
type CLOBTrade struct {
	ID                  int64           `db:"id" json:"id"`
	TradeID             uint64          `db:"trade_id" json:"trade_id"`
	MarketID            uint64          `db:"market_id" json:"market_id"`
	Buyer               string          `db:"buyer" json:"buyer"`
	BuyerSubAccountID   uint64          `db:"buyer_sub_account_id" json:"buyer_sub_account_id"`
	Seller              string          `db:"seller" json:"seller"`
	SellerSubAccountID  uint64          `db:"seller_sub_account_id" json:"seller_sub_account_id"`
	BuyerOrderID        uint64          `db:"buyer_order_id" json:"buyer_order_id"`
	SellerOrderID       uint64          `db:"seller_order_id" json:"seller_order_id"`
	Price               decimal.Decimal `db:"price" json:"price"`
	Quantity            decimal.Decimal `db:"quantity" json:"quantity"`
	TradeValue          decimal.Decimal `db:"trade_value" json:"trade_value"`
	BuyerFee            decimal.Decimal `db:"buyer_fee" json:"buyer_fee"`
	SellerFee           decimal.Decimal `db:"seller_fee" json:"seller_fee"`
	IsBuyerTaker        bool            `db:"is_buyer_taker" json:"is_buyer_taker"`
	IsBuyerLiquidation  bool            `db:"is_buyer_liquidation" json:"is_buyer_liquidation"`
	IsSellerLiquidation bool            `db:"is_seller_liquidation" json:"is_seller_liquidation"`
	ExecutedAt          time.Time       `db:"executed_at" json:"executed_at"`
	BlockHeight         int64           `db:"block_height" json:"block_height"`
	TxHash              string          `db:"tx_hash" json:"tx_hash"`
}

// CLOBPosition represents a CLOB perpetual position
type CLOBPosition struct {
	ID                  int64             `db:"id" json:"id"`
	PositionID          uint64            `db:"position_id" json:"position_id"`
	MarketID            uint64            `db:"market_id" json:"market_id"`
	Owner               string            `db:"owner" json:"owner"`
	SubAccountID        uint64            `db:"sub_account_id" json:"sub_account_id"`
	Side                CLOBPositionSide  `db:"side" json:"side"`
	Size                decimal.Decimal   `db:"size" json:"size"`
	Notional            decimal.Decimal   `db:"notional" json:"notional"`
	EntryPrice          decimal.Decimal   `db:"entry_price" json:"entry_price"`
	MarkPrice           decimal.Decimal   `db:"mark_price" json:"mark_price"`
	LiquidationPrice    decimal.Decimal   `db:"liquidation_price" json:"liquidation_price"`
	Margin              decimal.Decimal   `db:"margin" json:"margin"`
	MarginRatio         decimal.Decimal   `db:"margin_ratio" json:"margin_ratio"`
	UnrealizedPnL       decimal.Decimal   `db:"unrealized_pnl" json:"unrealized_pnl"`
	RealizedPnL         decimal.Decimal   `db:"realized_pnl" json:"realized_pnl"`
	CumulativeFunding   decimal.Decimal   `db:"cumulative_funding" json:"cumulative_funding"`
	LastFundingPayment  decimal.Decimal   `db:"last_funding_payment" json:"last_funding_payment"`
	LastFundingTime     *time.Time        `db:"last_funding_time" json:"last_funding_time,omitempty"`
	OpenedAt            time.Time         `db:"opened_at" json:"opened_at"`
	UpdatedAt           time.Time         `db:"updated_at" json:"updated_at"`
	ClosedAt            *time.Time        `db:"closed_at" json:"closed_at,omitempty"`
	ClosePrice          *decimal.Decimal  `db:"close_price" json:"close_price,omitempty"`
	IsLiquidated        bool              `db:"is_liquidated" json:"is_liquidated"`
	BlockHeight         int64             `db:"block_height" json:"block_height"`
	TxHash              string            `db:"tx_hash" json:"tx_hash"`
}

// CLOBOrderBookSnapshot represents an order book snapshot
type CLOBOrderBookSnapshot struct {
	ID              int64           `db:"id" json:"id"`
	MarketID        uint64          `db:"market_id" json:"market_id"`
	Bids            json.RawMessage `db:"bids" json:"bids"`
	Asks            json.RawMessage `db:"asks" json:"asks"`
	BestBid         *decimal.Decimal `db:"best_bid" json:"best_bid,omitempty"`
	BestAsk         *decimal.Decimal `db:"best_ask" json:"best_ask,omitempty"`
	MidPrice        *decimal.Decimal `db:"mid_price" json:"mid_price,omitempty"`
	Spread          *decimal.Decimal `db:"spread" json:"spread,omitempty"`
	TotalBidVolume  decimal.Decimal `db:"total_bid_volume" json:"total_bid_volume"`
	TotalAskVolume  decimal.Decimal `db:"total_ask_volume" json:"total_ask_volume"`
	SnapshotTime    time.Time       `db:"snapshot_time" json:"snapshot_time"`
	BlockHeight     int64           `db:"block_height" json:"block_height"`
}

// CLOBFundingRate represents funding rate data
type CLOBFundingRate struct {
	ID              int64           `db:"id" json:"id"`
	MarketID        uint64          `db:"market_id" json:"market_id"`
	FundingRate     decimal.Decimal `db:"funding_rate" json:"funding_rate"`
	PremiumRate     decimal.Decimal `db:"premium_rate" json:"premium_rate"`
	MarkPrice       decimal.Decimal `db:"mark_price" json:"mark_price"`
	IndexPrice      decimal.Decimal `db:"index_price" json:"index_price"`
	Timestamp       time.Time       `db:"timestamp" json:"timestamp"`
	NextFundingTime time.Time       `db:"next_funding_time" json:"next_funding_time"`
	BlockHeight     int64           `db:"block_height" json:"block_height"`
}

// CLOBLiquidation represents a liquidation event
type CLOBLiquidation struct {
	ID                      int64            `db:"id" json:"id"`
	LiquidationID           uint64           `db:"liquidation_id" json:"liquidation_id"`
	MarketID                uint64           `db:"market_id" json:"market_id"`
	PositionID              uint64           `db:"position_id" json:"position_id"`
	Owner                   string           `db:"owner" json:"owner"`
	SubAccountID            uint64           `db:"sub_account_id" json:"sub_account_id"`
	Liquidator              *string          `db:"liquidator" json:"liquidator,omitempty"`
	Side                    CLOBPositionSide `db:"side" json:"side"`
	Size                    decimal.Decimal  `db:"size" json:"size"`
	Price                   decimal.Decimal  `db:"price" json:"price"`
	LiquidationFee          decimal.Decimal  `db:"liquidation_fee" json:"liquidation_fee"`
	InsuranceFundContribution decimal.Decimal `db:"insurance_fund_contribution" json:"insurance_fund_contribution"`
	IsADL                   bool             `db:"is_adl" json:"is_adl"`
	LiquidatedAt            time.Time        `db:"liquidated_at" json:"liquidated_at"`
	BlockHeight             int64            `db:"block_height" json:"block_height"`
	TxHash                  string           `db:"tx_hash" json:"tx_hash"`
}

// CLOBMarketStats represents aggregated market statistics
type CLOBMarketStats struct {
	ID                   int64            `db:"id" json:"id"`
	MarketID             uint64           `db:"market_id" json:"market_id"`
	Volume24h            decimal.Decimal  `db:"volume_24h" json:"volume_24h"`
	Trades24h            int64            `db:"trades_24h" json:"trades_24h"`
	High24h              *decimal.Decimal `db:"high_24h" json:"high_24h,omitempty"`
	Low24h               *decimal.Decimal `db:"low_24h" json:"low_24h,omitempty"`
	OpenInterest         decimal.Decimal  `db:"open_interest" json:"open_interest"`
	OpenInterestNotional decimal.Decimal  `db:"open_interest_notional" json:"open_interest_notional"`
	LastPrice            *decimal.Decimal `db:"last_price" json:"last_price,omitempty"`
	LastTradeTime        *time.Time       `db:"last_trade_time" json:"last_trade_time,omitempty"`
	UpdatedAt            time.Time        `db:"updated_at" json:"updated_at"`
}

// OrderBookLevel represents a single level in the order book
type OrderBookLevel struct {
	Price    decimal.Decimal `json:"price"`
	Quantity decimal.Decimal `json:"quantity"`
	Orders   int             `json:"orders"`
}

// WebSocket update types for CLOB
type WSCLOBOrderUpdate struct {
	Action    string     `json:"action"`
	Order     *CLOBOrder `json:"order"`
	Timestamp time.Time  `json:"timestamp"`
}

type WSCLOBTradeUpdate struct {
	Trade     *CLOBTrade `json:"trade"`
	Timestamp time.Time  `json:"timestamp"`
}

type WSCLOBPositionUpdate struct {
	Action    string        `json:"action"`
	Position  *CLOBPosition `json:"position"`
	Timestamp time.Time     `json:"timestamp"`
}

type WSCLOBOrderBookUpdate struct {
	MarketID  uint64           `json:"market_id"`
	Bids      []OrderBookLevel `json:"bids"`
	Asks      []OrderBookLevel `json:"asks"`
	Timestamp time.Time        `json:"timestamp"`
}