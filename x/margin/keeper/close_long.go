package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CloseLong(ctx sdk.Context, msg *types.MsgClose) (*types.MTP, sdk.Int, error) {
	mtp, err := k.GetMTP(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Amm pool Id
	poolId := mtp.AmmPoolId

	// Get pool from pool Id
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
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

	if k.hooks != nil {
		k.hooks.AfterMarginPositionClosed(ctx, ammPool, pool)
	}

	return &mtp, repayAmount, nil
}
