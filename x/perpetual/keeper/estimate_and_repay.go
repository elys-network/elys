package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// EstimateAndRepay ammPool has to be pointer because RemoveFromPoolBalance (in Repay) updates pool assets
// Important to send pointer mtp and pool
func (k Keeper) EstimateAndRepay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool, baseCurrency string, closingRatio sdk.Dec) (math.Int, error) {

	if closingRatio.LTE(math.LegacyZeroDec()) || closingRatio.GT(math.LegacyOneDec()) {
		return math.Int{}, fmt.Errorf("invalid closing ratio (%s)", closingRatio.String())
	}

	repayAmount, payingLiabilities, err := k.CalcRepayAmount(ctx, mtp, ammPool, closingRatio)
	if err != nil {
		return math.ZeroInt(), err
	}
	returnAmount, err := k.CalcReturnAmount(*mtp, repayAmount, closingRatio)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// Note: Long settlement is done in trading asset. And short settlement in usdc in Repay function
	if err = k.Repay(ctx, mtp, pool, ammPool, returnAmount, payingLiabilities, closingRatio, baseCurrency); err != nil {
		return sdk.ZeroInt(), err
	}

	return repayAmount, nil
}

// CalcRepayAmount repay amount is in custody asset for liabilities with closing ratio
func (k Keeper) CalcRepayAmount(ctx sdk.Context, mtp *types.MTP, ammPool *ammtypes.Pool, closingRatio sdk.Dec) (repayAmount math.Int, payingLiabilities math.Int, err error) {
	// init repay amount
	// For long this will be in trading asset (custody asset is trading asset)
	// For short this will be in USDC (custody asset is USDC)
	repayAmount = sdk.ZeroInt()

	// mtp.BorrowInterestUnpaidLiability is 0 because settled in SettleInterest so no need to add
	// For short this will be in trading asset
	// For long this will be in base currency
	payingLiabilities = mtp.Liabilities.ToLegacyDec().Mul(closingRatio).TruncateInt()

	if mtp.Position == types.Position_LONG {
		liabilitiesWithClosingRatio := sdk.NewCoin(mtp.LiabilitiesAsset, payingLiabilities)
		repayAmount, _, err = k.EstimateSwapGivenOut(ctx, liabilitiesWithClosingRatio, mtp.CustodyAsset, *ammPool)
		if err != nil {
			return math.ZeroInt(), math.ZeroInt(), err
		}
	}
	if mtp.Position == types.Position_SHORT {
		// if position is short, repay in custody asset which is base currency
		liabilitiesWithClosingRatio := sdk.NewCoin(mtp.LiabilitiesAsset, payingLiabilities)
		repayAmount, _, err = k.EstimateSwapGivenOut(ctx, liabilitiesWithClosingRatio, mtp.CustodyAsset, *ammPool)
		if err != nil {
			return math.ZeroInt(), sdk.ZeroInt(), err
		}
	}

	return repayAmount, payingLiabilities, nil

}

// need to make sure unpaid liability interest is paid
func (k Keeper) CalcReturnAmount(mtp types.MTP, repayAmount math.Int, closingRatio sdk.Dec) (returnAmount math.Int, err error) {
	// closingAmount is what user is trying to close
	// For long, mtp.Custody is trading asset, unit of repay amount here is custody asset
	// For short mtp.Custody is base currency, unit of repay amount here is custody asset
	closingAmount := mtp.Custody.ToLegacyDec().Mul(closingRatio).TruncateInt()

	if closingAmount.LT(repayAmount) {
		// this case would mean bot liquidation failed as custody amount fall too low after interest was paid
		returnAmount = sdk.ZeroInt()
	} else {
		// can afford both
		returnAmount = closingAmount.Sub(repayAmount)
	}
	return returnAmount, nil
}
