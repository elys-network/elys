package types

const (
	// Event types
	EventDeposit          = "deposit"
	EventWithdraw         = "withdraw"
	EventUpdateParams     = "update_params"
	EventPlaceLimitOrder  = "place_limit_order"
	EventPlaceMarketOrder = "place_market_order"
	EventCreateMarket     = "create_perpetual_market"
	EventLiquidation      = "liquidation"
	EventOrderExecuted    = "order_executed"
	EventPositionOpened   = "position_opened"
	EventPositionClosed   = "position_closed"
	EventPositionModified = "position_modified"
	EventFundingPayment   = "funding_payment"

	// Attributes
	AttributeAmount        = "amount"
	AttributeAuthority     = "authority"
	AttributeSender        = "sender"
	AttributeMarketId      = "market_id"
	AttributeOrderId       = "order_id"
	AttributeOrderType     = "order_type"
	AttributePrice         = "price"
	AttributeQuantity      = "quantity"
	AttributePositionId    = "position_id"
	AttributeLiquidator    = "liquidator"
	AttributeOwner         = "owner"
	AttributeReward        = "reward"
	AttributePnL           = "pnl"
	AttributeFundingAmount = "funding_amount"
	AttributeTradePrice    = "trade_price"
	AttributeTradeQuantity = "trade_quantity"
	AttributeSide          = "side"
	AttributeIsLiquidation = "is_liquidation"
)
