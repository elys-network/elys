package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// SettleFundingFeeCollection handles funding fee collection
func (k Keeper) SettleFundingFeeCollection(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) error {
	// get funding rate
	fundingRate := k.GetFundingRate(ctx, mtp.LastInterestCalcBlock, mtp.AmmPoolId)

	// if funding rate is zero, return
	if fundingRate.IsZero() {
		return nil
	}

	// if funding rate is negative and mtp position is long or funding rate is positive and mtp position is short, return
	if (fundingRate.IsNegative() && mtp.Position == types.Position_LONG) || (fundingRate.IsPositive() && mtp.Position == types.Position_SHORT) {
		return nil
	}

	// Calculate the take amount in custody asset
	takeAmountCustodyAmount := types.CalcTakeAmount(mtp.Custody, mtp.CustodyAsset, fundingRate)

	// Build the take amount coin
	takeAmountCustody := sdk.NewCoin(mtp.CustodyAsset, takeAmountCustodyAmount)

	// Swap the take amount to collateral asset
	takeAmountCollateralAmount, err := k.EstimateSwap(ctx, takeAmountCustody, mtp.CollateralAsset, ammPool)
	if err != nil {
		return err
	}

	// Get funding fee collection address
	fundingFeeCollectionAddress := k.GetFundingFeeCollectionAddress(ctx)

	// Transfer take amount in collateral asset to funding fee collection address
	_, err = k.TakeFundPayment(ctx, takeAmountCollateralAmount, mtp.CollateralAsset, sdk.OneDec(), fundingFeeCollectionAddress, &ammPool)
	if err != nil {
		return err
	}

	// update mtp custody
	mtp.Custody = mtp.Custody.Sub(takeAmountCustodyAmount)

	// add payment to total funding fee paid in collateral asset
	mtp.FundingFeePaidCollateral = mtp.FundingFeePaidCollateral.Add(takeAmountCollateralAmount)

	// add payment to total funding fee paid in custody asset
	mtp.FundingFeePaidCustody = mtp.FundingFeePaidCustody.Add(takeAmountCustodyAmount)

	// emit event
	if !takeAmountCollateralAmount.IsZero() {
		k.EmitFundingFeePayment(ctx, mtp, takeAmountCustody.Amount, mtp.CollateralAsset, types.EventIncrementalPayFund)
	}

	// update pool custody balance
	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, takeAmountCustody.Amount, false, mtp.Position)
	if err != nil {
		return err
	}

	// update accounted balance for custody side
	err = pool.UpdateBalance(ctx, mtp.CustodyAsset, takeAmountCustody.Amount, false, mtp.Position)
	if err != nil {
		return err
	}

	// update accounted balance for collateral side
	err = pool.UpdateBalance(ctx, mtp.CollateralAsset, takeAmountCollateralAmount, false, mtp.Position)
	if err != nil {
		return err
	}

	// apply changes to mtp object
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	// apply changes to pool object
	k.SetPool(ctx, *pool)

	// update mtp health
	_, err = k.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return err
	}

	return nil
}
