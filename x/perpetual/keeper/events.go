package keeper

import (
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) EmitFundPayment(ctx sdk.Context, mtp *types.MTP, takeAmount math.Int, takeAsset string, paymentType string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(paymentType,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("payment_amount", takeAmount.String()),
		sdk.NewAttribute("payment_asset", takeAsset),
	))
}

func (k Keeper) EmitForceClose(ctx sdk.Context, mtp *types.MTP, repayAmount math.Int, closer string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventForceClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collaterals", mtp.Collateral.String()),
		sdk.NewAttribute("custodies", mtp.Custody.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", fmt.Sprintf("%s", mtp.Leverage)),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("borrow_interest_paid_collaterals", mtp.BorrowInterestPaidCollateral.String()),
		sdk.NewAttribute("borrow_interest_paid_custodies", mtp.BorrowInterestPaidCustody.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
		sdk.NewAttribute("closer", closer),
	))
}

func (k Keeper) EmitFundingFeePayment(ctx sdk.Context, mtp *types.MTP, takeAmount math.Int, takeAsset string, paymentType string) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(paymentType,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("payment_amount", takeAmount.String()),
		sdk.NewAttribute("payment_asset", takeAsset),
	))
}
