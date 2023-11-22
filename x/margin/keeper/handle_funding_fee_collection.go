package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

// HandleFundingFeeCollection handles funding fee collection
func (k Keeper) HandleFundingFeeCollection(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, collateralAsset string, custodyAsset string) error {
	// get funding rate
	fundingRate := pool.FundingRate

	// if funding rate is zero, return
	if fundingRate.IsZero() {
		return nil
	}

	// if funding rate is negative and mtp position is long or funding rate is positive and mtp position is short, return
	if (fundingRate.IsNegative() && mtp.Position == types.Position_LONG) || (fundingRate.IsPositive() && mtp.Position == types.Position_SHORT) {
		return nil
	}

	// get indexes
	collateralIndex, custodyIndex := types.GetMTPAssetIndex(mtp, collateralAsset, custodyAsset)

	// Calculate the take amount in custody asset
	takeAmountCustody := types.CalcTakeAmount(mtp.Custodies[custodyIndex], custodyAsset, fundingRate)

	// Swap the take amount to collateral asset
	takeAmountCollateralAmount, err := k.EstimateSwap(ctx, takeAmountCustody, collateralAsset, ammPool)
	if err != nil {
		return err
	}

	// Create the take amount coin
	takeAmountCollateral := sdk.NewCoin(collateralAsset, takeAmountCollateralAmount)

	// Get funding fee collection address
	fundingFeeCollectionAddress := k.GetFundingFeeCollectionAddress(ctx)

	// Transfer take amount in collateral asset to funding fee collection address
	_, err = k.TakeFundPayment(ctx, takeAmountCollateral.Amount, collateralAsset, sdk.OneDec(), fundingFeeCollectionAddress, &ammPool)
	if err != nil {
		return err
	}

	// update mtp custody
	mtp.Custodies[custodyIndex] = mtp.Custodies[custodyIndex].Sub(takeAmountCustody)

	// add payment to total funding fee paid in collateral asset
	mtp.FundingFeePaidCollaterals[collateralIndex] = mtp.FundingFeePaidCollaterals[collateralIndex].Add(takeAmountCollateral)

	// add payment to total funding fee paid in custody asset
	mtp.FundingFeePaidCustodies[custodyIndex] = mtp.FundingFeePaidCustodies[custodyIndex].Add(takeAmountCustody)

	// emit event
	if !takeAmountCollateral.IsZero() {
		k.EmitFundingFeePayment(ctx, mtp, takeAmountCustody.Amount, collateralAsset, types.EventIncrementalPayFund)
	}

	// update pool custody balance
	err = pool.UpdateCustody(ctx, custodyAsset, takeAmountCustody.Amount, false, mtp.Position)
	if err != nil {
		return err
	}

	// update accounted balance for custody side
	err = pool.UpdateBalance(ctx, custodyAsset, takeAmountCustody.Amount, false, mtp.Position)
	if err != nil {
		return err
	}

	// update accounted balance for collateral side
	err = pool.UpdateBalance(ctx, collateralAsset, takeAmountCollateral.Amount, false, mtp.Position)
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
	_, err = k.UpdateMTPHealth(ctx, *mtp, ammPool, custodyAsset)
	if err != nil {
		return err
	}

	return nil
}
