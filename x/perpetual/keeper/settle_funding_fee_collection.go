package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// SettleFunding handles funding fee collection and distribution
func (k Keeper) SettleFunding(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) error {

	err := k.SettleFundingFeeCollection(ctx, mtp, pool, ammPool, baseCurrency)
	if err != nil {
		return err
	}

	err = k.SettleFundingFeeDistribution(ctx, mtp, pool, ammPool, baseCurrency)
	if err != nil {
		return err
	}

	mtp.LastFundingCalcBlock = uint64(ctx.BlockHeight())
	mtp.LastFundingCalcTime = uint64(ctx.BlockTime().Unix())

	// apply changes to mtp object
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) SettleFundingFeeCollection(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) error {
	// get funding rate
	_, longRate, shortRate := k.GetFundingRate(ctx, mtp.LastFundingCalcBlock, mtp.AmmPoolId)

	var takeAmountCustodyAmount math.Int
	if mtp.Position == types.Position_LONG {
		takeAmountCustodyAmount = types.CalcTakeAmount(mtp.Custody, longRate)
	} else {
		takeAmountCustodyAmount = types.CalcTakeAmount(mtp.Custody, shortRate)
	}
	// Calculate the take amount in custody asset
	if !takeAmountCustodyAmount.IsPositive() {
		return nil
	}

	takeAmountCustody := sdk.NewCoin(mtp.CustodyAsset, takeAmountCustodyAmount)

	// Swap the take amount to collateral asset
	takeAmountCollateralAmount, _, err := k.EstimateSwap(ctx, takeAmountCustody, mtp.CollateralAsset, ammPool)
	if err != nil {
		return err
	}

	// increase fees collected
	err = pool.UpdateFeesCollected(ctx, mtp.CollateralAsset, takeAmountCollateralAmount, true)
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

	// apply changes to mtp object
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	// apply changes to pool object
	k.SetPool(ctx, *pool)

	return nil
}
