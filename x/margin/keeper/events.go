package keeper

import (
	"fmt"
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
		sdk.NewAttribute("collaterals", fmt.Sprintf("%s", mtp.Collaterals)),
		sdk.NewAttribute("custodies", fmt.Sprintf("%s", mtp.Custodies)),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", fmt.Sprintf("%s", mtp.Leverages)),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("borrow_interest_paid_collaterals", fmt.Sprintf("%s", mtp.BorrowInterestPaidCollaterals)),
		sdk.NewAttribute("borrow_interest_paid_custodies", fmt.Sprintf("%s", mtp.BorrowInterestPaidCustodies)),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
		sdk.NewAttribute("closer", closer),
	))
}

func (k Keeper) EmitFundingFeePayment(ctx sdk.Context, mtp *types.MTP, takeAmount sdk.Int, takeAsset string, paymentType string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(paymentType,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("payment_amount", takeAmount.String()),
		sdk.NewAttribute("payment_asset", takeAsset),
	))
}
