package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// EstimateAndRepay ammPool has to be pointer because RemoveFromPoolBalance (in Repay) updates pool assets
// Important to send pointer mtp and pool
func (k Keeper) EstimateAndRepay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool, closingRatio math.LegacyDec) (math.Int, math.Int, error) {

	if closingRatio.LTE(math.LegacyZeroDec()) || closingRatio.GT(math.LegacyOneDec()) {
		return math.Int{}, math.Int{}, fmt.Errorf("invalid closing ratio (%s)", closingRatio.String())
	}

	repayAmount, payingLiabilities, _, slippageAmount, weightBreakingFee, repayOracleAmount, perpetualFees, takerFees, err := k.CalcRepayAmount(ctx, mtp, ammPool, closingRatio)
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), err
	}
	k.CalculateAndEmitPerpetualFeesEvent(ctx, ammPool.PoolParams.UseOracle, sdk.NewCoin(mtp.CustodyAsset, repayAmount), sdk.NewCoin(mtp.LiabilitiesAsset, payingLiabilities), slippageAmount, weightBreakingFee, perpetualFees, takerFees, repayOracleAmount, false)

	returnAmount, err := k.CalcReturnAmount(*mtp, repayAmount, closingRatio)
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), err
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return math.Int{}, math.Int{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// Note: Long settlement is done in trading asset. And short settlement in usdc in Repay function
	if err = k.Repay(ctx, mtp, pool, ammPool, returnAmount, payingLiabilities, closingRatio, baseCurrency); err != nil {
		return math.ZeroInt(), math.ZeroInt(), err
	}

	return repayAmount, returnAmount, nil
}

// CalcRepayAmount repay amount is in custody asset for liabilities with closing ratio
func (k Keeper) CalcRepayAmount(ctx sdk.Context, mtp *types.MTP, ammPool *ammtypes.Pool, closingRatio math.LegacyDec) (repayAmount, payingLiabilities math.Int, slippage, slippageAmount, weightBreakingFee, repayOracleAmount osmomath.BigDec, perpetualFees, takerFees math.LegacyDec, err error) {
	// init repay amount
	// For long this will be in trading asset (custody asset is trading asset)
	// For short this will be in USDC (custody asset is USDC)
	repayAmount = math.ZeroInt()

	// mtp.BorrowInterestUnpaidLiability is 0 because settled in SettleInterest so no need to add
	// For short this will be in trading asset
	// For long this will be in base currency
	payingLiabilities = closingRatio.MulInt(mtp.Liabilities).TruncateInt()

	if mtp.Position == types.Position_LONG {
		liabilitiesWithClosingRatio := sdk.NewCoin(mtp.LiabilitiesAsset, payingLiabilities)
		repayAmount, slippage, slippageAmount, weightBreakingFee, repayOracleAmount, perpetualFees, takerFees, err = k.EstimateSwapGivenOut(ctx, liabilitiesWithClosingRatio, mtp.CustodyAsset, *ammPool, mtp.Address)
		if err != nil {
			return math.ZeroInt(), math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
	}
	if mtp.Position == types.Position_SHORT {
		// if position is short, repay in custody asset which is base currency
		liabilitiesWithClosingRatio := sdk.NewCoin(mtp.LiabilitiesAsset, payingLiabilities)
		repayAmount, slippage, slippageAmount, weightBreakingFee, repayOracleAmount, perpetualFees, takerFees, err = k.EstimateSwapGivenOut(ctx, liabilitiesWithClosingRatio, mtp.CustodyAsset, *ammPool, mtp.Address)
		if err != nil {
			return math.ZeroInt(), math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
	}
	k.TrackSlippageAndWeightBreakingSlippage(ctx, ammPool, slippage, weightBreakingFee, repayAmount, mtp.CustodyAsset)

	return

}

// need to make sure unpaid liability interest is paid
func (k Keeper) CalcReturnAmount(mtp types.MTP, repayAmount math.Int, closingRatio math.LegacyDec) (returnAmount math.Int, err error) {
	// closingAmount is what user is trying to close
	// For long, mtp.Custody is trading asset, unit of repay amount here is custody asset
	// For short mtp.Custody is base currency, unit of repay amount here is custody asset
	closingAmount := closingRatio.MulInt(mtp.Custody).TruncateInt()

	if closingAmount.LT(repayAmount) {
		// this case would mean bot liquidation failed as custody amount fall too low after interest was paid
		returnAmount = math.ZeroInt()
	} else {
		// can afford both
		returnAmount = closingAmount.Sub(repayAmount)
	}
	return returnAmount, nil
}
