package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) ClosePosition(ctx sdk.Context, msg *types.MsgClose) (types.MTP, math.Int, osmomath.BigDec, math.Int, math.Int, math.Int, math.Int, math.Int, bool, bool, error) {
	// Retrieve MTP
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.GetMTP(ctx, creator, msg.Id)
	if err != nil {
		return types.MTP{}, math.ZeroInt(), osmomath.ZeroBigDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, err
	}

	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return mtp, math.ZeroInt(), osmomath.ZeroBigDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", mtp.AmmPoolId))
	}

	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return mtp, math.ZeroInt(), osmomath.ZeroBigDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, err
	}

	// this also handles edge case where bot is unable to close position in time.
	repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, err := k.MTPTriggerChecksAndUpdates(ctx, &mtp, &pool, &ammPool)
	if err != nil {
		return types.MTP{}, math.ZeroInt(), osmomath.ZeroBigDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, err
	}

	if forceClosed {
		return mtp, repayAmt, osmomath.OneBigDec(), returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, nil
	}

	// Should be declared after SettleMTPBorrowInterestUnpaidLiability and settling funding
	closingRatio := osmomath.BigDecFromSDKInt(msg.Amount).Quo(mtp.GetBigDecCustody())
	if mtp.Position == types.Position_SHORT {
		closingRatio = osmomath.BigDecFromSDKInt(msg.Amount).Quo(mtp.GetBigDecLiabilities())
	}
	if closingRatio.GT(osmomath.OneBigDec()) {
		closingRatio = osmomath.OneBigDec()
	}

	// Estimate swap and repay
	repayAmt, returnAmt, err = k.EstimateAndRepay(ctx, &mtp, &pool, &ammPool, closingRatio)
	if err != nil {
		return mtp, math.ZeroInt(), osmomath.ZeroBigDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), allInterestsPaid, forceClosed, err
	}

	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		params := k.GetParams(ctx)
		err = k.hooks.AfterPerpetualPositionClosed(ctx, ammPool, pool, creator, params.EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return mtp, math.Int{}, osmomath.BigDec{}, math.Int{}, math.Int{}, math.Int{}, math.Int{}, math.Int{}, allInterestsPaid, forceClosed, err
		}
	}

	return mtp, repayAmt, closingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, nil
}
