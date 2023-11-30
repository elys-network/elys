package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CloseLong(ctx sdk.Context, msg *types.MsgClose, baseCurrency string) (*types.MTP, sdk.Int, error) {
	// Retrieve MTP
	mtp, err := k.CloseLongChecker.GetMTP(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Retrieve Pool
	pool, found := k.CloseLongChecker.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return nil, sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	// Retrieve AmmPool
	ammPool, err := k.CloseLongChecker.GetAmmPool(ctx, mtp.AmmPoolId, mtp.CustodyAsset)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Handle Borrow Interest if within epoch position
	if err := k.CloseLongChecker.HandleBorrowInterest(ctx, &mtp, &pool, ammPool); err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Take out custody
	err = k.CloseLongChecker.TakeOutCustody(ctx, mtp, &pool, msg.Amount)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Estimate swap and repay
	repayAmt, err := k.CloseLongChecker.EstimateAndRepay(ctx, mtp, pool, ammPool, msg.Amount, baseCurrency)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Hooks after margin position closed
	if k.hooks != nil {
		k.hooks.AfterMarginPositionClosed(ctx, ammPool, pool)
	}

	return &mtp, repayAmt, nil
}
