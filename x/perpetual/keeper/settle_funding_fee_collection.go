package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) FundingFeeCollection(ctx sdk.Context, mtp *types.MTP, pool *types.Pool) (bool, math.Int, error) {

	fullFundingFeePayment := true
	var takeAmountCustodyAmount math.Int
	// get funding rate
	longRate, shortRate := k.GetFundingRate(ctx, mtp.LastFundingCalcBlock, mtp.LastFundingCalcTime, mtp.AmmPoolId)

	if mtp.Position == types.Position_LONG {
		takeAmountCustodyAmount = types.CalcTakeAmount(mtp.Custody, longRate)
		if !takeAmountCustodyAmount.IsPositive() {
			return true, math.ZeroInt(), nil
		}

		if takeAmountCustodyAmount.GT(mtp.Custody) {
			fullFundingFeePayment = false
			takeAmountCustodyAmount = mtp.Custody
		}

		// increase fees collected
		err := pool.UpdateFeesCollected(mtp.CustodyAsset, takeAmountCustodyAmount, true)
		if err != nil {
			return fullFundingFeePayment, math.ZeroInt(), err
		}

		// update mtp custody
		mtp.Custody = mtp.Custody.Sub(takeAmountCustodyAmount)

		// add payment to total funding fee paid in custody asset
		mtp.FundingFeePaidCustody = mtp.FundingFeePaidCustody.Add(takeAmountCustodyAmount)

		// update pool custody balance
		err = pool.UpdateCustody(mtp.CustodyAsset, takeAmountCustodyAmount, false, mtp.Position)
		if err != nil {
			return fullFundingFeePayment, math.ZeroInt(), err
		}
	} else {
		takeAmountLiabilityAmount := types.CalcTakeAmount(mtp.Liabilities, shortRate)
		if !takeAmountLiabilityAmount.IsPositive() {
			return true, math.ZeroInt(), nil
		}

		// increase fees collected
		// Note: fees is collected in liabilities asset
		err := pool.UpdateFeesCollected(mtp.LiabilitiesAsset, takeAmountLiabilityAmount, true)
		if err != nil {
			return fullFundingFeePayment, math.ZeroInt(), err
		}

		_, tradingAssetPriceBaseDenomRatio, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
		if err != nil {
			return fullFundingFeePayment, math.ZeroInt(), err
		}

		// should be done in custody
		// short -> usdc
		// long -> custody
		// For short, takeAmountLiabilityAmount is in trading asset, need to convert to custody asset which is in usdc
		takeAmountCustodyAmount = osmomath.BigDecFromSDKInt(takeAmountLiabilityAmount).Mul(tradingAssetPriceBaseDenomRatio).Dec().TruncateInt()

		if takeAmountCustodyAmount.GT(mtp.Custody) {
			fullFundingFeePayment = false
			takeAmountCustodyAmount = mtp.Custody
		}

		// update mtp custody
		mtp.Custody = mtp.Custody.Sub(takeAmountCustodyAmount)
		// add payment to total funding fee paid in custody asset
		mtp.FundingFeePaidCustody = mtp.FundingFeePaidCustody.Add(takeAmountCustodyAmount)

		// update pool custody balance
		err = pool.UpdateCustody(mtp.CustodyAsset, takeAmountCustodyAmount, false, mtp.Position)
		if err != nil {
			return fullFundingFeePayment, math.ZeroInt(), err
		}
	}

	return fullFundingFeePayment, takeAmountCustodyAmount, nil
}
