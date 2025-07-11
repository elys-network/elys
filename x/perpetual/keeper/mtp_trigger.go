package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

// MTPTriggerChecksAndUpdates Should be run whenever there is a state change of MTP. Runs in following order:
// 1. Settle funding fee, if unable to settle whole, close the position
// 2. Settle interest payments, in unable to pay whole, close the position
// 3. Update position health and check if above minimum threshold, if equal or lower close the position
func (k Keeper) MTPTriggerChecksAndUpdates(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool) (repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt math.Int, allInterestsPaid, forceClosed bool, perpetualFeesCoins types.PerpetualFees, closingPrice math.LegacyDec, err error) {

	allInterestsPaid = true
	forceClosed = false
	fundingFeeFullyPaid := true
	interestFullyPaid := true

	// Update interests
	err = k.UpdateMTPBorrowInterestUnpaidLiability(ctx, mtp)
	if err != nil {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, sdkerrors.Wrap(err, "error updating borrow interest unpaid liability")
	}

	// Pay funding fee
	fundingFeeFullyPaid, fundingFeeAmt, fundingAmtDistributed, err = k.SettleFunding(ctx, mtp, pool)
	if err != nil {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, sdkerrors.Wrap(err, "error handling funding fee")
	}

	// Unable to pay funding fee, close the position
	if !fundingFeeFullyPaid {
		allInterestsPaid = false
		forceClosed = true
		repayAmt, returnAmt, perpetualFeesCoins, closingPrice, err = k.ForceClose(ctx, mtp, pool, ammPool)
		if err != nil {
			return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, sdkerrors.Wrap(err, "error executing force close")
		}

		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, math.ZeroInt(), math.ZeroInt(), allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, nil
	}

	// Pay interests
	interestAmt, insuranceAmt, interestFullyPaid, err = k.SettleMTPBorrowInterestUnpaidLiability(ctx, mtp, pool, ammPool)
	if err != nil {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, sdkerrors.Wrap(err, "error handling borrow interest payment")
	}

	// Unable to pay interests, close the position
	if !interestFullyPaid {
		allInterestsPaid = false
		forceClosed = true
		repayAmt, returnAmt, perpetualFeesCoins, closingPrice, err = k.ForceClose(ctx, mtp, pool, ammPool)
		if err != nil {
			return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, sdkerrors.Wrap(err, "error executing force close")
		}
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, nil

	}

	mtp.MtpHealth, err = k.GetMTPHealth(ctx, *mtp)
	if err != nil {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, sdkerrors.Wrap(err, "error updating mtp health")
	}

	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, err
	}

	k.SetPool(ctx, *pool)

	safetyFactor := k.GetSafetyFactor(ctx)
	// Position is unhealthy, close the position
	if mtp.MtpHealth.LTE(safetyFactor) {
		forceClosed = true
		repayAmt, returnAmt, perpetualFeesCoins, closingPrice, err = k.ForceClose(ctx, mtp, pool, ammPool)
		if err != nil {
			return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, sdkerrors.Wrap(err, "error executing force close")
		}
	}

	return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, nil
}
