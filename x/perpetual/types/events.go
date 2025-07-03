package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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
	EventPerpetualFees             = "perpetual_fees"
)

const (
	AttributeKeyPerpFee           = "perp_fee"
	AttributeKeySlippage          = "slippage"
	AttributeKeyWeightBreakingFee = "weight_breaking_fee"
	AttributeTakerFees            = "taker_fees"
)

func EmitPerpetualFeesEvent(ctx sdk.Context, perpFee, slippage, weightBreakingFee, takerFees string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		NewPerpFeeEvent(perpFee, slippage, weightBreakingFee, takerFees),
	})
}

func NewPerpFeeEvent(perpFee, slippage, weightBreakingFee, takerFees string) sdk.Event {
	return sdk.NewEvent(
		EventPerpetualFees,
		sdk.NewAttribute("value", "USD"),
		sdk.NewAttribute(AttributeKeyPerpFee, perpFee),
		sdk.NewAttribute(AttributeKeySlippage, slippage),
		sdk.NewAttribute(AttributeKeyWeightBreakingFee, weightBreakingFee),
		sdk.NewAttribute(AttributeTakerFees, takerFees),
	)
}
