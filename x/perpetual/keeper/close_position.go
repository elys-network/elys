package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) ClosePosition(ctx sdk.Context, msg *types.MsgClose) (types.MTP, math.Int, math.LegacyDec, math.Int, math.Int, math.Int, math.Int, bool, bool, error) {
	// Retrieve MTP
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.GetMTP(ctx, creator, msg.Id)
	if err != nil {
		return types.MTP{}, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, err
	}

	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return mtp, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", mtp.AmmPoolId))
	}

	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return mtp, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, err
	}

	// this also handles edge case where bot is unable to close position in time.
	repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, err := k.MTPTriggerChecksAndUpdates(ctx, &mtp, &pool, &ammPool)
	if err != nil {
		return types.MTP{}, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, err
	}

	if forceClosed {
		return mtp, repayAmt, math.LegacyOneDec(), returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, nil
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
	repayAmt, returnAmt, err = k.EstimateAndRepay(ctx, &mtp, &pool, &ammPool, closingRatio)
	if err != nil {
		return mtp, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), allInterestsPaid, forceClosed, err
	}

	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		params := k.GetParams(ctx)
		err = k.hooks.AfterPerpetualPositionClosed(ctx, ammPool, pool, creator, params.EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return mtp, math.Int{}, math.LegacyDec{}, math.Int{}, math.Int{}, math.Int{}, math.Int{}, allInterestsPaid, forceClosed, err
		}
	}

	return mtp, repayAmt, closingRatio, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, nil
}
