package keeper

import (
	"errors"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// MTPTriggerChecksAndUpdates Should be run whenever there is a state change of MTP. Runs in following order:
// 1. Settle funding fee, if unable to settle whole, close the position
// 2. Settle interest payments, in unable to pay whole, close the position
// 3. Update position health and check if above minimum threshold, if equal or lower close the position
func (k Keeper) MTPTriggerChecksAndUpdates(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool) (repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt math.Int, allInterestsPaid, forceClosed bool, err error) {

	allInterestsPaid = true
	forceClosed = false
	fundingFeeFullyPaid := true
	interestFullyPaid := true

	// Update interests
	k.UpdateMTPBorrowInterestUnpaidLiability(ctx, mtp)

	// Pay funding fee
	fundingFeeFullyPaid, fundingFeeAmt, fundingAmtDistributed, err = k.SettleFunding(ctx, mtp, pool)
	if err != nil {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, sdkerrors.Wrap(err, "error handling funding fee")
	}

	// Unable to pay funding fee, close the position
	if !fundingFeeFullyPaid {
		allInterestsPaid = false
		forceClosed = true
		repayAmt, returnAmt, err = k.ForceClose(ctx, mtp, pool, ammPool)
		if err != nil {
			return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, sdkerrors.Wrap(err, "error executing force close")
		}

		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, math.ZeroInt(), math.ZeroInt(), allInterestsPaid, forceClosed, nil
	}

	// Pay interests
	interestAmt, insuranceAmt, interestFullyPaid, err = k.SettleMTPBorrowInterestUnpaidLiability(ctx, mtp, pool, ammPool)
	if err != nil {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, sdkerrors.Wrap(err, "error handling borrow interest payment")
	}

	// Unable to pay interests, close the position
	if !interestFullyPaid {
		allInterestsPaid = false
		forceClosed = true
		repayAmt, returnAmt, err = k.ForceClose(ctx, mtp, pool, ammPool)
		if err != nil {
			return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, sdkerrors.Wrap(err, "error executing force close")
		}
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, nil
	}

	baseCurrencyEntry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, errors.New("unable to find base currency entry")
	}

	baseCurrency := baseCurrencyEntry.Denom

	h, err := k.GetMTPHealth(ctx, *mtp, *ammPool, baseCurrency)
	if err != nil {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, sdkerrors.Wrap(err, "error updating mtp health")
	}
	mtp.MtpHealth = h.Dec()

	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, err
	}

	k.SetPool(ctx, *pool)

	safetyFactor := k.GetSafetyFactor(ctx)
	// Position is unhealthy, close the position
	if mtp.MtpHealth.LTE(safetyFactor) {
		forceClosed = true
		repayAmt, returnAmt, err = k.ForceClose(ctx, mtp, pool, ammPool)
		if err != nil {
			return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, sdkerrors.Wrap(err, "error executing force close")
		}
	}

	return repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, nil
}
