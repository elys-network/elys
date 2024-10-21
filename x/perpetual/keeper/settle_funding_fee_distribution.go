package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) FundingFeeDistribution(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) error {
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
			return fmt.Errorf("totalCustodyLong in FundingFeeDistribution cannot be zero")
		}
		fundingFeeShare = mtp.Custody.ToLegacyDec().Quo(totalCustodyLong.ToLegacyDec())
		totalFund = short

		// if funding fee share is zero, skip mtp
		if fundingFeeShare.IsZero() || totalFund.IsZero() {
			return nil
		}

		// calculate funding fee amount
		fundingFeeAmount := totalFund.Mul(fundingFeeShare)

		// update mtp custody
		mtp.Custody = mtp.Custody.Add(fundingFeeAmount.TruncateInt())

		// decrease fees collected
		err := pool.UpdateFeesCollected(mtp.CustodyAsset, fundingFeeAmount.TruncateInt(), false)
		if err != nil {
			return err
		}

		// update pool custody balance
		err = pool.UpdateCustody(mtp.CustodyAsset, fundingFeeAmount.TruncateInt(), true, mtp.Position)
		if err != nil {
			return err
		}

		// add payment to total funding fee paid in custody asset
		mtp.FundingFeeReceivedCustody = mtp.FundingFeeReceivedCustody.Add(fundingFeeAmount.TruncateInt())
	} else {
		// Ensure liabilitiesShort is not zero to avoid division by zero
		if totalLiabilitiesShort.IsZero() {
			return types.ErrAmountTooLow
		}
		fundingFeeShare = mtp.Liabilities.ToLegacyDec().Quo(totalLiabilitiesShort.ToLegacyDec())
		totalFund = long

		// if funding fee share is zero, skip mtp
		if fundingFeeShare.IsZero() || totalFund.IsZero() {
			return nil
		}

		// calculate funding fee amount
		fundingFeeAmount := totalFund.Mul(fundingFeeShare).TruncateInt()

		// adding case for fundingFeeAmount being smaller tha 10^-18
		if fundingFeeAmount.IsZero() {
			return nil
		}
		// decrease fees collected
		err := pool.UpdateFeesCollected(mtp.LiabilitiesAsset, fundingFeeAmount, false)
		if err != nil {
			return err
		}

		custodyAmt, _, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.LiabilitiesAsset, fundingFeeAmount), mtp.CustodyAsset, ammPool)
		if err != nil {
			return err
		}
		// update mtp Custody
		mtp.Custody = mtp.Custody.Add(custodyAmt)

		// update pool liability balance
		err = pool.UpdateCustody(mtp.CustodyAsset, custodyAmt, true, mtp.Position)
		if err != nil {
			return err
		}

		// add payment to total funding fee paid in custody asset
		mtp.FundingFeeReceivedCustody = mtp.FundingFeeReceivedCustody.Add(custodyAmt)
	}

	return nil
}
