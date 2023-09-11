package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	mtp, err := k.GetMTP(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, err
	}

	var closedMtp *types.MTP
	var repayAmount sdk.Int
	switch mtp.Position {
	case types.Position_LONG:
		closedMtp, repayAmount, err = k.CloseLong(ctx, msg)
		if err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		closedMtp, repayAmount, err = k.CloseShort(ctx, msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, mtp.Position.String())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(closedMtp.Id), 10)),
		sdk.NewAttribute("position", closedMtp.Position.String()),
		sdk.NewAttribute("address", closedMtp.Address),
		sdk.NewAttribute("collateral_asset", closedMtp.CollateralAsset),
		sdk.NewAttribute("collateral_amount", closedMtp.CollateralAmount.String()),
		sdk.NewAttribute("custody_asset", closedMtp.CustodyAsset),
		sdk.NewAttribute("custody_amount", closedMtp.CustodyAmount.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", closedMtp.Leverage.String()),
		sdk.NewAttribute("liabilities", closedMtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollateral.String()),
		sdk.NewAttribute("interest_paid_custody", mtp.InterestPaidCustody.String()),
		sdk.NewAttribute("interest_unpaid_collateral", closedMtp.InterestUnpaidCollateral.String()),
		sdk.NewAttribute("health", closedMtp.MtpHealth.String()),
	))

	return &types.MsgCloseResponse{}, nil
}
