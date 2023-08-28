package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k msgServer) Close(goCtx context.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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

func (k msgServer) CloseLong(ctx sdk.Context, msg *types.MsgClose) (*types.MTP, sdk.Int, error) {
	mtp, err := k.GetMTP(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}
	// Get pool Ids which can support borrowing asset
	poolIds := k.amm.GetAllPoolIdsWithDenom(ctx, mtp.CustodyAsset)
	if len(poolIds) < 1 {
		return nil, sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}

	// Assume choose the first pool
	poolId := poolIds[0]

	// Get pool from pool Id
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}

	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return nil, sdk.ZeroInt(), sdkerrors.Wrap(types.ErrPoolDoesNotExist, mtp.CustodyAsset)
	}

	epochLength := k.GetEpochLength(ctx)
	epochPosition := GetEpochPosition(ctx, epochLength)
	if epochPosition > 0 {
		interestPayment := CalcMTPInterestLiabilities(&mtp, pool.InterestRate, epochPosition, epochLength)
		finalInterestPayment := k.HandleInterestPayment(ctx, interestPayment, &mtp, &pool, ammPool)

		err = pool.UpdateBlockInterest(ctx, mtp.CollateralAsset, finalInterestPayment, true)
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}

		mtp.MtpHealth, err = k.UpdateMTPHealth(ctx, mtp, ammPool)
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}
	}

	err = k.TakeOutCustody(ctx, mtp, &pool)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	cutodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, mtp.CustodyAmount)
	repayAmount, err := k.EstimateSwap(ctx, cutodyAmtTokenIn, mtp.CollateralAsset, ammPool)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	err = k.Repay(ctx, &mtp, &pool, ammPool, repayAmount, false)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	return &mtp, repayAmount, nil
}
