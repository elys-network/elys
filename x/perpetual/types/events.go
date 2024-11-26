package types

const (
	EventOpen                     = "perpetual/mtp_open"
	EventClose                    = "perpetual/mtp_close"
	EventForceClose               = "perpetual/mtp_force_close"
	EventIncrementalPayFund       = "perpetual/incremental_pay_fund"
	EventRepayFund                = "perpetual/repay_fund"
	EventClosePositions           = "perpetual/close_positions"
	EventClosePositionsTakeProfit = "perpetual/close_positions_take_profit"
	EventClosePositionsUnhealthy  = "perpetual/close_positions_unhealthy"
	EventClosePositionsStopLoss   = "perpetual/close_positions_stop_loss"
)
