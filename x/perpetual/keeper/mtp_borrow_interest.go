package keeper

import (
	"errors"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetBorrowInterestAmount(ctx sdk.Context, mtp *types.MTP) math.Int {

	err := mtp.UpdateMTPTakeProfitBorrowFactor()
	if err != nil {
		panic(err)
	}

	// This already gives a floor tested value for interest rate
	borrowInterestRate := k.GetBorrowInterestRate(ctx, mtp.LastInterestCalcBlock, mtp.LastInterestCalcTime, mtp.AmmPoolId, mtp.GetBigDecTakeProfitBorrowFactor())
	totalLiability := mtp.Liabilities.Add(mtp.BorrowInterestUnpaidLiability)
	borrowInterestPayment := osmomath.BigDecFromSDKInt(totalLiability).Mul(borrowInterestRate).Dec().TruncateInt()
	return borrowInterestPayment
}

func (k Keeper) UpdateMTPBorrowInterestUnpaidLiability(ctx sdk.Context, mtp *types.MTP) {
	borrowInterestPaymentInt := k.GetBorrowInterestAmount(ctx, mtp)
	mtp.BorrowInterestUnpaidLiability = mtp.BorrowInterestUnpaidLiability.Add(borrowInterestPaymentInt)
	mtp.LastInterestCalcBlock = uint64(ctx.BlockHeight())
	mtp.LastInterestCalcTime = uint64(ctx.BlockTime().Unix())
}

// SettleMTPBorrowInterestUnpaidLiability  This does not update BorrowInterestUnpaidLiability, it should be done through UpdateMTPBorrowInterestUnpaidLiability beforehand
func (k Keeper) SettleMTPBorrowInterestUnpaidLiability(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool) (math.Int, math.Int, bool, error) {

	// adding case for mtp.BorrowInterestUnpaidLiability being smaller tha 10^-18, this happens when position is small so liabilities are small, and hardly any time has spend, so interests are small, so it leads to 0 value
	if mtp.BorrowInterestUnpaidLiability.IsZero() {
		// nothing to pay back
		return math.ZeroInt(), math.ZeroInt(), true, nil
	}

	_, tradingAssetPriceDenomRatio, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), true, err
	}

	borrowInterestPaymentInCustody, err := mtp.GetBorrowInterestAmountAsCustodyAsset(tradingAssetPriceDenomRatio)
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), true, err
	}

	// here we are paying the interests so unpaid borrow interest reset to 0
	mtp.BorrowInterestUnpaidLiability = math.ZeroInt()

	fullyPaid := true
	if borrowInterestPaymentInCustody.GT(mtp.Custody) {
		unpaidInterestCustody := borrowInterestPaymentInCustody.Sub(mtp.Custody)
		unpaidInterestLiabilities := math.ZeroInt()
		if mtp.Position == types.Position_LONG {
			// custody is in trading asset, liabilities needs to be in usdc,
			unpaidInterestLiabilities = osmomath.BigDecFromSDKInt(unpaidInterestCustody).Mul(tradingAssetPriceDenomRatio).Dec().TruncateInt()
		} else {
			// custody is in usdc, liabilities needs to be in trading asset,
			if tradingAssetPriceDenomRatio.IsZero() {
				return math.ZeroInt(), math.ZeroInt(), false, errors.New("trading asset price is zero in SettleMTPBorrowInterestUnpaidLiability")
			}
			borrowInterestPaymentInCustody = osmomath.BigDecFromSDKInt(unpaidInterestCustody).Quo(tradingAssetPriceDenomRatio).Dec().TruncateInt()
		}
		mtp.BorrowInterestUnpaidLiability = unpaidInterestLiabilities

		// Since not enough custody left, we can only pay this much, position health goes to 0
		borrowInterestPaymentInCustody = mtp.Custody
		fullyPaid = false

		insuranceBalance := k.bankKeeper.GetBalance(ctx, pool.GetInsuranceAccount(), mtp.CustodyAsset)
		if insuranceBalance.Amount.GTE(unpaidInterestCustody) {
			coin := sdk.NewCoin(mtp.CustodyAsset, unpaidInterestCustody)
			err = k.SendToAmmPool(ctx, pool.GetInsuranceAccount(), ammPool, sdk.NewCoins(coin))
			if err != nil {
				return math.ZeroInt(), math.ZeroInt(), false, err
			}
			ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventPaidFromInsuranceFund,
				sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(mtp.Id), 10)),
				sdk.NewAttribute("owner", mtp.Address),
				sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(mtp.AmmPoolId), 10)),
				sdk.NewAttribute("position", mtp.Position.String()),
				sdk.NewAttribute("amount", coin.String()),
			))
		} else {
			ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventInsufficientInsuranceFund,
				sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(mtp.Id), 10)),
				sdk.NewAttribute("owner", mtp.Address),
				sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(mtp.AmmPoolId), 10)),
				sdk.NewAttribute("position", mtp.Position.String()),
				sdk.NewAttribute("amount", insuranceBalance.String()),
				sdk.NewAttribute("denom", mtp.CustodyAsset),
				sdk.NewAttribute("unpaid_interest_custody", unpaidInterestCustody.String()),
			))
		}
	}

	mtp.BorrowInterestPaidCustody = mtp.BorrowInterestPaidCustody.Add(borrowInterestPaymentInCustody)

	// deduct borrow interest payment from custody amount
	// This will go to zero if borrowInterestPaymentInCustody.GT(mtp.Custody) true
	mtp.Custody = mtp.Custody.Sub(borrowInterestPaymentInCustody)

	insuranceAmount := math.ZeroInt()
	// if full interest is paid then only collect insurance fund
	if fullyPaid {
		insuranceAmount, err = k.CollectInsuranceFund(ctx, borrowInterestPaymentInCustody, mtp.CustodyAsset, ammPool, *pool)
		if err != nil {
			return math.ZeroInt(), math.ZeroInt(), fullyPaid, err
		}
	}

	err = pool.UpdateCustody(mtp.CustodyAsset, borrowInterestPaymentInCustody, false, mtp.Position)
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), fullyPaid, err
	}

	return borrowInterestPaymentInCustody, insuranceAmount, fullyPaid, nil

}
