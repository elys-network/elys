package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	mtp, err := k.GetMTP(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, err
	}

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	var closedMtp *types.MTP
	var repayAmount sdk.Int
	switch mtp.Position {
	case types.Position_LONG:
		closedMtp, repayAmount, err = k.CloseLong(ctx, msg, baseCurrency)
		if err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		closedMtp, repayAmount, err = k.CloseShort(ctx, msg, baseCurrency)
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
		sdk.NewAttribute("collateral", closedMtp.Collateral.String()),
		sdk.NewAttribute("custody", closedMtp.Custody.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", fmt.Sprintf("%s", closedMtp.Leverage)),
		sdk.NewAttribute("liabilities", closedMtp.Liabilities.String()),
		sdk.NewAttribute("borrow_interest_paid_collateral", mtp.BorrowInterestPaidCollateral.String()),
		sdk.NewAttribute("borrow_interest_paid_custody", mtp.BorrowInterestPaidCustody.String()),
		sdk.NewAttribute("borrow_interest_unpaid_collateral", closedMtp.BorrowInterestUnpaidCollateral.String()),
		sdk.NewAttribute("health", closedMtp.MtpHealth.String()),
	))

	return &types.MsgCloseResponse{}, nil
}
