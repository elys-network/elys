package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CloseShort(ctx sdk.Context, msg *types.MsgClose) (*types.MTP, sdk.Int, error) {
	// Retrieve MTP
	mtp, err := k.CloseShortChecker.GetMTP(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Retrieve Pool
	pool, found := k.CloseShortChecker.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return nil, sdk.ZeroInt(), sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	repayAmount := sdk.ZeroInt()
	for _, custodyAsset := range mtp.CustodyAssets {
		// Retrieve AmmPool
		ammPool, err := k.CloseShortChecker.GetAmmPool(ctx, mtp.AmmPoolId, custodyAsset)
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}

		for _, collateralAsset := range mtp.CollateralAssets {

			// Handle Interest if within epoch position
			if err := k.CloseShortChecker.HandleInterest(ctx, &mtp, &pool, ammPool, collateralAsset, custodyAsset); err != nil {
				return nil, sdk.ZeroInt(), err
			}
		}

		// Take out custody
		err = k.CloseShortChecker.TakeOutCustody(ctx, mtp, &pool, custodyAsset)
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}

		for _, collateralAsset := range mtp.CollateralAssets {
			// Estimate swap and repay
			repayAmt, err := k.CloseShortChecker.EstimateAndRepay(ctx, mtp, pool, ammPool, collateralAsset, custodyAsset)
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
