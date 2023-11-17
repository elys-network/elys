package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CloseLong(ctx sdk.Context, msg *types.MsgClose) (*types.MTP, sdk.Int, error) {
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

	repayAmount := sdk.ZeroInt()
	for _, custody := range mtp.Custodies {
		custodyAsset := custody.Denom
		// Retrieve AmmPool
		ammPool, err := k.CloseLongChecker.GetAmmPool(ctx, mtp.AmmPoolId, custodyAsset)
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}

		for _, collateral := range mtp.Collaterals {
			collateralAsset := collateral.Denom
			// Handle Borrow Interest if within epoch position
			if err := k.CloseLongChecker.HandleBorrowInterest(ctx, &mtp, &pool, ammPool, collateralAsset, custodyAsset); err != nil {
				return nil, sdk.ZeroInt(), err
			}
		}

		// Take out custody
		err = k.CloseLongChecker.TakeOutCustody(ctx, mtp, &pool, custodyAsset)
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}

		for _, collateral := range mtp.Collaterals {
			collateralAsset := collateral.Denom
			// Estimate swap and repay
			repayAmt, err := k.CloseLongChecker.EstimateAndRepay(ctx, mtp, pool, ammPool, collateralAsset, custodyAsset)
			if err != nil {
				return nil, sdk.ZeroInt(), err
			}

			repayAmount = repayAmount.Add(repayAmt)
		}

		// Hooks after margin position closed
		if k.hooks != nil {
			k.hooks.AfterMarginPositionClosed(ctx, ammPool, pool)
		}
	}

	return &mtp, repayAmount, nil
}
