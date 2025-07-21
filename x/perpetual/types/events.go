package types

const (
	EventOpen                      = "perpetual_mtp_open"
	EventUpdateStopLoss            = "perpetual_mtp_update_stop_loss"
	EventUpdateTakeProfitPrice     = "perpetual_mtp_update_take_profit_price"
	EventOpenConsolidate           = "perpetual_mtp_open_consolidate"
	EventClose                     = "perpetual_mtp_close"
	EventForceClosed               = "perpetual_mtp_force_closed"
	EventPaidFromInsuranceFund     = "perpetual_mtp_paid_from_insurance_fund"
	EventInsufficientInsuranceFund = "perpetual_mtp_insufficient_insurance_fund"
	EventAddCollateral             = "perpetual_mtp_add_collateral"
	EventADLClose                  = "perpetual_adl_close"
	EventClosePositions            = "perpetual_close_positions"
)

const (
	AttributeKeyPerpFee           = "perp_fee"
	AttributeKeySlippage          = "slippage"
	AttributeKeyWeightBreakingFee = "weight_breaking_fee"
	AttributeTakerFees            = "taker_fees"
)
