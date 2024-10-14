package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) SettleFundingFeeDistribution(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) error {
	// account custody from long position
	totalCustodyLong := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsLong {
		totalCustodyLong = totalCustodyLong.Add(asset.Custody)
	}

	// account custody from short position
	totalLiabilitiesShort := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsShort {
		totalLiabilitiesShort = totalLiabilitiesShort.Add(asset.Liabilities)
	}

	// Total fund collected should be
	long, short := k.GetFundingDistributionValue(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId)
	var totalFund sdk.Dec
	// calc funding fee share
	var fundingFeeShare sdk.Dec
	if mtp.Position == types.Position_LONG {
		// Ensure liabilitiesLong is not zero to avoid division by zero
		if totalCustodyLong.IsZero() {
			return types.ErrAmountTooLow
		}
		fundingFeeShare = sdk.NewDecFromInt(mtp.Custody).Quo(sdk.NewDecFromInt(totalCustodyLong))
		totalFund = short
	} else {
		// Ensure liabilitiesShort is not zero to avoid division by zero
		if totalLiabilitiesShort.IsZero() {
			return types.ErrAmountTooLow
		}
		fundingFeeShare = sdk.NewDecFromInt(mtp.Liabilities).Quo(sdk.NewDecFromInt(totalLiabilitiesShort))
		totalFund = long
	}

	// if funding fee share is zero, skip mtp
	if fundingFeeShare.IsZero() {
		return nil
	}

	// calculate funding fee amount
	fundingFeeAmount := totalFund.Mul(fundingFeeShare)

	// update mtp custody
	// TODO: Check for short position
	mtp.Custody = mtp.Custody.Add(fundingFeeAmount.TruncateInt())

	// decrease fees collected
	err := pool.UpdateFeesCollected(ctx, mtp.CustodyAsset, fundingFeeAmount.TruncateInt(), false)
	if err != nil {
		return err
	}

	// update pool custody balance
	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, fundingFeeAmount.TruncateInt(), true, mtp.Position)
	if err != nil {
		return err
	}

	// add payment to total funding fee paid in custody asset
	mtp.FundingFeeReceivedCustody = mtp.FundingFeeReceivedCustody.Add(fundingFeeAmount.TruncateInt())

	return nil
}
