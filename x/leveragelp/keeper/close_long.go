package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
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
	for _, custodyAsset := range mtp.CustodyAssets {
		// Retrieve AmmPool
		ammPool, err := k.CloseLongChecker.GetAmmPool(ctx, mtp.AmmPoolId, custodyAsset)
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}

		for _, collateralAsset := range mtp.CollateralAssets {
			// Handle Interest if within epoch position
			if err := k.CloseLongChecker.HandleInterest(ctx, &mtp, &pool, ammPool, collateralAsset, custodyAsset); err != nil {
				return nil, sdk.ZeroInt(), err
			}
		}

		// Take out custody
		err = k.CloseLongChecker.TakeOutCustody(ctx, mtp, &pool, custodyAsset)
		if err != nil {
			return nil, sdk.ZeroInt(), err
		}

		for _, collateralAsset := range mtp.CollateralAssets {
			// Estimate swap and repay
			repayAmt, err := k.CloseLongChecker.EstimateAndRepay(ctx, mtp, pool, ammPool, collateralAsset, custodyAsset)
			if err != nil {
				return nil, sdk.ZeroInt(), err
			}

			repayAmount = repayAmount.Add(repayAmt)
		}

		// Hooks after leveragelp position closed
		if k.hooks != nil {
			k.hooks.AfterLeveragelpPositionClosed(ctx, ammPool, pool)
		}
	}

	return &mtp, repayAmount, nil
}
