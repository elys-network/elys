package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) EmitFundPayment(ctx sdk.Context, mtp *types.MTP, takeAmount sdk.Int, takeAsset string, paymentType string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(paymentType,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("payment_amount", takeAmount.String()),
		sdk.NewAttribute("payment_asset", takeAsset),
	))
}

func (k Keeper) EmitForceClose(ctx sdk.Context, mtp *types.MTP, repayAmount sdk.Int, closer string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventForceClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("address", mtp.Address),
		// sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		// sdk.NewAttribute("collateral_amount", mtp.CollateralAmount.String()),
		// sdk.NewAttribute("custody_asset", mtp.CustodyAsset),
		// sdk.NewAttribute("custody_amount", mtp.CustodyAmount.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		// sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		// sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollateral.String()),
		// sdk.NewAttribute("interest_paid_custody", mtp.InterestPaidCustody.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
		sdk.NewAttribute("closer", closer),
	))
}
