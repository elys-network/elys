package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) ClosePosition(ctx sdk.Context, msg *types.MsgClose, baseCurrency string) (*types.MTP, math.Int, error) {
	// Retrieve MTP
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.GetMTP(ctx, creator, msg.Id)
	if err != nil {
		return nil, math.ZeroInt(), err
	}

	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return nil, math.ZeroInt(), err
	}

	// This needs to be updated here to check user doesn't send more than required amount
	k.UpdateMTPBorrowInterestUnpaidLiability(ctx, &mtp)
	// Retrieve Pool
	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return nil, math.ZeroInt(), errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", mtp.AmmPoolId))
	}

	// Handle Borrow Interest if within epoch position SettleMTPBorrowInterestUnpaidLiability settles interest using mtp.Custody, mtp.Custody gets reduced
	if _, err = k.SettleMTPBorrowInterestUnpaidLiability(ctx, &mtp, &pool, ammPool); err != nil {
		return nil, math.ZeroInt(), err
	}

	err = k.SettleFunding(ctx, &mtp, &pool, ammPool)
	if err != nil {
		return nil, math.ZeroInt(), errorsmod.Wrapf(err, "error handling funding fee")
	}

	// Should be declared after SettleMTPBorrowInterestUnpaidLiability and settling funding
	closingRatio := msg.Amount.ToLegacyDec().Quo(mtp.Custody.ToLegacyDec())
	if mtp.Position == types.Position_SHORT {
		closingRatio = msg.Amount.ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec())
	}
	if closingRatio.GT(math.LegacyOneDec()) {
		closingRatio = math.LegacyOneDec()
	}

	// Estimate swap and repay
	repayAmt, err := k.EstimateAndRepay(ctx, &mtp, &pool, &ammPool, baseCurrency, closingRatio)
	if err != nil {
		return nil, math.ZeroInt(), err
	}

	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		params := k.GetParams(ctx)
		err = k.hooks.AfterPerpetualPositionClosed(ctx, ammPool, pool, creator, params.EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return nil, math.Int{}, err
		}
	}

	return &mtp, repayAmt, nil
}
