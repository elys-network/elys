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

	collateralIndex := len(mtp.Collaterals) - 1
	custodyIndex := len(mtp.Custodies) - 1
	mtpPosIndex := len(mtp.Leverages) - 1

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(closedMtp.Id), 10)),
		sdk.NewAttribute("position", closedMtp.Position.String()),
		sdk.NewAttribute("address", closedMtp.Address),
		sdk.NewAttribute("collateral", closedMtp.Collaterals[collateralIndex].String()),
		sdk.NewAttribute("custody", closedMtp.Custodies[custodyIndex].String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", closedMtp.Leverages[mtpPosIndex].String()),
		sdk.NewAttribute("liabilities", closedMtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollaterals[collateralIndex].String()),
		sdk.NewAttribute("interest_paid_custody", mtp.InterestPaidCustodies[custodyIndex].String()),
		sdk.NewAttribute("interest_unpaid_collateral", closedMtp.InterestUnpaidCollaterals[collateralIndex].String()),
		sdk.NewAttribute("health", closedMtp.MtpHealth.String()),
	))

	return &types.MsgCloseResponse{}, nil
}
