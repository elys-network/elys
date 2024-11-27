package types

const (
	EventOpen                 = "perpetual/mtp_open"
	EventClose                = "perpetual/mtp_close"
	EventForceCloseUnhealthy  = "perpetual/mtp_force_close_unhealthy"
	EventForceCloseStopLoss   = "perpetual/mtp_force_close_stopLoss"
	EventForceCloseTakeprofit = "perpetual/mtp_force_close_takeprofit"
	EventIncrementalPayFund   = "perpetual/incremental_pay_fund"
	EventRepayFund            = "perpetual/repay_fund"
	EventClosePositions       = "perpetual/close_positions"
)
