package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) ClosePosition(ctx sdk.Context, msg *types.MsgClose, baseCurrency string) (*types.MTP, math.Int, error) {
	// Retrieve MTP
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.ClosePositionChecker.GetMTP(ctx, creator, msg.Id)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	if msg.Amount.GT(mtp.Custody) || msg.Amount.IsNegative() {
		return nil, sdk.ZeroInt(), types.ErrInvalidCloseSize
	}

	// Retrieve Pool
	pool, found := k.ClosePositionChecker.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return nil, sdk.ZeroInt(), errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	// Retrieve AmmPool
	ammPool, err := k.ClosePositionChecker.GetAmmPool(ctx, mtp.AmmPoolId, mtp.CustodyAsset)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Handle Borrow Interest if within epoch position
	if _, err := k.ClosePositionChecker.SettleBorrowInterest(ctx, &mtp, &pool, ammPool); err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Take out custody
	err = k.ClosePositionChecker.TakeOutCustody(ctx, mtp, &pool, msg.Amount)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Estimate swap and repay
	repayAmt, err := k.ClosePositionChecker.EstimateAndRepay(ctx, mtp, pool, ammPool, msg.Amount, baseCurrency)
	if err != nil {
		return nil, sdk.ZeroInt(), err
	}

	// Hooks after perpetual position closed
	if k.hooks != nil {
		k.hooks.AfterPerpetualPositionClosed(ctx, ammPool, pool, creator)
	}

	return &mtp, repayAmt, nil
}
