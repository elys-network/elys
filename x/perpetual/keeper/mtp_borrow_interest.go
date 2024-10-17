package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) GetBorrowInterestAmount(ctx sdk.Context, mtp *types.MTP) math.Int {

	err := k.UpdateMTPTakeProfitBorrowFactor(ctx, mtp)
	if err != nil {
		panic(err)
	}

	borrowInterestRate := k.GetBorrowInterestRate(ctx, mtp.LastInterestCalcBlock, mtp.LastInterestCalcTime, mtp.AmmPoolId, mtp.TakeProfitBorrowFactor)
	borrowInterestRateWithProfitFactor := borrowInterestRate.Mul(mtp.TakeProfitBorrowFactor)
	totalLiability := mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability)
	borrowInterestPayment := totalLiability.ToLegacyDec().Mul(borrowInterestRateWithProfitFactor).TruncateInt()
	minBorrowInterestAmount := k.GetParams(ctx).MinBorrowInterestAmount
	return sdk.MaxInt(borrowInterestPayment, minBorrowInterestAmount)
}

func (k Keeper) UpdateMTPBorrowInterestUnpaidLiability(ctx sdk.Context, mtp *types.MTP) {
	borrowInterestPaymentInt := k.GetBorrowInterestAmount(ctx, mtp)
	mtp.BorrowInterestUnpaidLiability = mtp.BorrowInterestUnpaidLiability.Add(borrowInterestPaymentInt)
	mtp.LastInterestCalcBlock = uint64(ctx.BlockHeight())
	mtp.LastInterestCalcTime = uint64(ctx.BlockTime().Unix())
	return
}

// SettleMTPBorrowInterestUnpaidLiability  This does not update BorrowInterestUnpaidLiability, it should be done through UpdateMTPBorrowInterestUnpaidLiability beforehand
func (k Keeper) SettleMTPBorrowInterestUnpaidLiability(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) (math.Int, error) {

	if mtp.BorrowInterestUnpaidLiability.IsZero() {
		return math.Int{}, fmt.Errorf("0 unpaid interest liability to settle")
	}

	liabilityInterestTokenOut := sdk.NewCoin(mtp.LiabilitiesAsset, mtp.BorrowInterestUnpaidLiability)
	borrowInterestPaymentInCustody, _, err := k.EstimateSwapGivenOut(ctx, liabilityInterestTokenOut, mtp.CustodyAsset, ammPool)
	if err != nil {
		return math.Int{}, errorsmod.Wrapf(err, "unable to swap BorrowInterestUnpaidLiability to custody asset (%s)", liabilityInterestTokenOut.String())
	}
	// here we are paying the interests so unpaid borrow interest reset to 0
	mtp.BorrowInterestUnpaidLiability = sdk.ZeroInt()

	// edge case, not enough custody to cover payment
	// TODO This should not happen, bot should close the position beforehand
	// TODO what if bot misses it, can we do anything about it?
	if borrowInterestPaymentInCustody.GT(mtp.Custody) {
		// TODO Do we need to keep this swap? as health will already be 0 as custody will be 0
		// TODO Doing this swap to update back mtp.BorrowInterestUnpaidLiability again as there aren't enough custody
		unpaidInterestCustody := borrowInterestPaymentInCustody.Sub(mtp.Custody)
		unpaidInterestCustodyTokenOut := sdk.NewCoin(mtp.CustodyAsset, unpaidInterestCustody)
		unpaidInterestLiabilities, _, err := k.EstimateSwapGivenOut(ctx, unpaidInterestCustodyTokenOut, mtp.LiabilitiesAsset, ammPool)
		if err != nil {
			return sdk.ZeroInt(), err
		}
		mtp.BorrowInterestUnpaidLiability = unpaidInterestLiabilities

		// Since not enough custody left, we can only pay this much, position health goes to 0
		borrowInterestPaymentInCustody = mtp.Custody
	}

	mtp.BorrowInterestPaidCustody = mtp.BorrowInterestPaidCustody.Add(borrowInterestPaymentInCustody)

	// deduct borrow interest payment from custody amount
	// This will go to zero if borrowInterestPaymentInCustody.GT(mtp.Custody) true
	mtp.Custody = mtp.Custody.Sub(borrowInterestPaymentInCustody)

	takePercentage := k.GetIncrementalBorrowInterestPaymentFundPercentage(ctx)
	fundAddr := k.GetIncrementalBorrowInterestPaymentFundAddress(ctx)
	takeAmount, err := k.TakeFundPayment(ctx, borrowInterestPaymentInCustody, mtp.CustodyAsset, takePercentage, fundAddr, &ammPool)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	actualBorrowInterestPaymentCustody := borrowInterestPaymentInCustody.Sub(takeAmount)

	if !takeAmount.IsZero() {
		k.EmitFundPayment(ctx, mtp, takeAmount, mtp.CustodyAsset, types.EventIncrementalPayFund)
	}

	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, borrowInterestPaymentInCustody, false, mtp.Position)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	return actualBorrowInterestPaymentCustody, nil

}
