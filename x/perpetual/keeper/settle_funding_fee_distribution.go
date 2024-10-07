package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) SettleFundingFeeDistribution(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) error {
	uusdc, found := k.assetProfileKeeper.GetEntry(ctx, "uusdc")
	if !found {
		return nil
	}

	// Calculate liabilities for long and short assets using the separate helper function
	liabilitiesLong, err := k.CalcTotalLiabilities(ctx, pool.PoolAssetsLong, pool.AmmPoolId, uusdc.Denom)
	if err != nil {
		return nil
	}

	liabilitiesShort, err := k.CalcTotalLiabilities(ctx, pool.PoolAssetsShort, pool.AmmPoolId, uusdc.Denom)
	if err != nil {
		return nil
	}

	// get funding fee collection address
	fundingFeeCollectionAddress := k.GetFundingFeeCollectionAddress(ctx)

	// get base currency balance
	balance := k.bankKeeper.GetBalance(ctx, fundingFeeCollectionAddress, baseCurrency)

	// if balance is zero, return
	if balance.IsZero() {
		return nil
	}

	// Total fund collected should be
	long, short := k.GetFundingDistributionValue(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId)
	var totalFund sdk.Dec
	// calc funding fee share
	var fundingFeeShare sdk.Dec
	if mtp.Position == types.Position_LONG {
		// Ensure liabilitiesLong is not zero to avoid division by zero
		if liabilitiesLong.IsZero() {
			return types.ErrAmountTooLow
		}
		fundingFeeShare = sdk.NewDecFromInt(mtp.Liabilities).Quo(sdk.NewDecFromInt(liabilitiesLong))
		totalFund = short
	} else {
		// Ensure liabilitiesShort is not zero to avoid division by zero
		if liabilitiesShort.IsZero() {
			return types.ErrAmountTooLow
		}
		fundingFeeShare = sdk.NewDecFromInt(mtp.Liabilities).Quo(sdk.NewDecFromInt(liabilitiesShort))
		totalFund = long
	}

	// if funding fee share is zero, skip mtp
	if fundingFeeShare.IsZero() {
		return nil
	}

	// calculate funding fee amount
	fundingFeeAmount := sdk.NewCoin(baseCurrency, totalFund.Mul(fundingFeeShare).TruncateInt())

	// update mtp custody
	mtp.Custody = mtp.Custody.Sub(fundingFeeAmount.Amount)

	// update pool custody balance
	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, fundingFeeAmount.Amount, false, mtp.Position)
	if err != nil {
		return err
	}

	// update accounted balance for custody side
	err = pool.UpdateBalance(ctx, mtp.CustodyAsset, fundingFeeAmount.Amount, false, mtp.Position)
	if err != nil {
		return err
	}

	// update received funding fee accounting buckets
	// Swap the take amount to collateral asset
	fundingFeeCollateralAmount, err := k.EstimateSwap(ctx, fundingFeeAmount, mtp.CollateralAsset, ammPool)
	if err != nil {
		return err
	}

	// add payment to total funding fee paid in collateral asset
	mtp.FundingFeeReceivedCollateral = mtp.FundingFeeReceivedCollateral.Add(fundingFeeCollateralAmount)
	// add payment to total funding fee paid in custody asset
	mtp.FundingFeeReceivedCustody = mtp.FundingFeeReceivedCustody.Add(fundingFeeAmount.Amount)

	return nil
}
