package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) GetBorrowInterestAmount(ctx sdk.Context, mtp *types.MTP) math.Int {

	err := mtp.UpdateMTPTakeProfitBorrowFactor()
	if err != nil {
		panic(err)
	}

	// This already gives a floor tested value for interest rate
	borrowInterestRate := k.GetBorrowInterestRate(ctx, mtp.LastInterestCalcBlock, mtp.LastInterestCalcTime, mtp.AmmPoolId, mtp.TakeProfitBorrowFactor)
	totalLiability := mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability)
	borrowInterestPayment := totalLiability.ToLegacyDec().Mul(borrowInterestRate).TruncateInt()
	return borrowInterestPayment
}

func (k Keeper) UpdateMTPBorrowInterestUnpaidLiability(ctx sdk.Context, mtp *types.MTP) {
	borrowInterestPaymentInt := k.GetBorrowInterestAmount(ctx, mtp)
	mtp.BorrowInterestUnpaidLiability = mtp.BorrowInterestUnpaidLiability.Add(borrowInterestPaymentInt)
	mtp.LastInterestCalcBlock = uint64(ctx.BlockHeight())
	mtp.LastInterestCalcTime = uint64(ctx.BlockTime().Unix())
}

// SettleMTPBorrowInterestUnpaidLiability  This does not update BorrowInterestUnpaidLiability, it should be done through UpdateMTPBorrowInterestUnpaidLiability beforehand
func (k Keeper) SettleMTPBorrowInterestUnpaidLiability(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) (math.Int, error) {

	// adding case for mtp.BorrowInterestUnpaidLiability being smaller tha 10^-18, this happens when position is small so liabilities are small, and hardly any time has spend, so interests are small, so it leads to 0 value
	if mtp.BorrowInterestUnpaidLiability.IsZero() {
		// nothing to pay back
		return math.ZeroInt(), nil
	}

	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return math.ZeroInt(), err
	}

	borrowInterestPaymentInCustody, err := mtp.GetBorrowInterestAmountAsCustodyAsset(tradingAssetPrice)
	if err != nil {
		return math.ZeroInt(), err
	}

	// here we are paying the interests so unpaid borrow interest reset to 0
	mtp.BorrowInterestUnpaidLiability = math.ZeroInt()

	// edge case, not enough custody to cover payment
	// TODO This should not happen, bot should close the position beforehand
	// TODO what if bot misses it, can we do anything about it?
	if borrowInterestPaymentInCustody.GT(mtp.Custody) {
		// TODO Do we need to keep this swap? as health will already be 0 as custody will be 0
		// TODO Doing this swap to update back mtp.BorrowInterestUnpaidLiability again as there aren't enough custody
		unpaidInterestCustody := borrowInterestPaymentInCustody.Sub(mtp.Custody)
		unpaidInterestLiabilities := math.ZeroInt()
		if mtp.Position == types.Position_LONG {
			// custody is in trading asset, liabilities needs to be in usdc,
			unpaidInterestLiabilities = unpaidInterestCustody.ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()
		} else {
			// custody is in usdc, liabilities needs to be in trading asset,
			borrowInterestPaymentInCustody = unpaidInterestCustody.ToLegacyDec().Quo(tradingAssetPrice).TruncateInt()
		}
		mtp.BorrowInterestUnpaidLiability = unpaidInterestLiabilities

		// Since not enough custody left, we can only pay this much, position health goes to 0
		borrowInterestPaymentInCustody = mtp.Custody
	}

	mtp.BorrowInterestPaidCustody = mtp.BorrowInterestPaidCustody.Add(borrowInterestPaymentInCustody)

	// deduct borrow interest payment from custody amount
	// This will go to zero if borrowInterestPaymentInCustody.GT(mtp.Custody) true
	mtp.Custody = mtp.Custody.Sub(borrowInterestPaymentInCustody)

	takeAmount, err := k.TakeFundPayment(ctx, borrowInterestPaymentInCustody, mtp.CustodyAsset, &ammPool)
	if err != nil {
		return math.ZeroInt(), err
	}

	if !takeAmount.IsZero() {
		k.EmitFundPayment(ctx, mtp, takeAmount, mtp.CustodyAsset, types.EventIncrementalPayFund)
	}

	err = pool.UpdateCustody(mtp.CustodyAsset, takeAmount, false, mtp.Position)
	if err != nil {
		return math.ZeroInt(), err
	}

	return borrowInterestPaymentInCustody, nil

}
