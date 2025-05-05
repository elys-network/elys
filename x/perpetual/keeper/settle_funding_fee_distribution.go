package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) FundingFeeDistribution(ctx sdk.Context, mtp *types.MTP, pool *types.Pool) (sdkmath.Int, error) {

	// Total fund collected should be
	longCollectedShare, shortCollectedShare := k.GetFundingDistributionValue(ctx, mtp.LastFundingCalcBlock, pool.AmmPoolId)
	amountDistributed := sdkmath.ZeroInt()
	if mtp.Position == types.Position_LONG {
		fundingFeeAmount := mtp.GetBigDecCustody().Mul(shortCollectedShare)
		if fundingFeeAmount.IsZero() {
			return amountDistributed, nil
		}

		amountDistributed = fundingFeeAmount.Dec().TruncateInt()
		// update mtp custody
		mtp.Custody = mtp.Custody.Add(fundingFeeAmount.Dec().TruncateInt())

		// decrease fees collected
		err := pool.UpdateFeesCollected(mtp.CustodyAsset, fundingFeeAmount.Dec().TruncateInt(), false)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}

		// update pool custody balance
		err = pool.UpdateCustody(mtp.CustodyAsset, fundingFeeAmount.Dec().TruncateInt(), true, mtp.Position)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}

		// add payment to total funding fee paid in custody asset
		mtp.FundingFeeReceivedCustody = mtp.FundingFeeReceivedCustody.Add(fundingFeeAmount.Dec().TruncateInt())
	} else {
		fundingFeeAmount := mtp.GetBigDecLiabilities().Mul(longCollectedShare).Dec().TruncateInt()

		// adding case for fundingFeeAmount being smaller tha 10^-18
		if fundingFeeAmount.IsZero() {
			return amountDistributed, nil
		}
		// decrease fees collected
		err := pool.UpdateFeesCollected(mtp.LiabilitiesAsset, fundingFeeAmount, false)
		if err != nil {
			return amountDistributed, err
		}

		_, tradingAssetPriceBaseDenomRatio, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
		if err != nil {
			return amountDistributed, err
		}

		// For short, fundingFeeAmount is in trading asset, need to convert to custody asset which is in usdc
		custodyAmt := osmomath.BigDecFromSDKInt(fundingFeeAmount).Mul(tradingAssetPriceBaseDenomRatio).Dec().TruncateInt()

		amountDistributed = custodyAmt
		// update mtp Custody
		mtp.Custody = mtp.Custody.Add(custodyAmt)

		// update pool liability balance
		err = pool.UpdateCustody(mtp.CustodyAsset, custodyAmt, true, mtp.Position)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}

		// add payment to total funding fee paid in custody asset
		mtp.FundingFeeReceivedCustody = mtp.FundingFeeReceivedCustody.Add(custodyAmt)
	}

	return amountDistributed, nil
}
