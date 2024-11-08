package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) FundingFeeCollection(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) error {
	// get funding rate
	longRate, shortRate := k.GetFundingRate(ctx, mtp.LastFundingCalcBlock, mtp.LastFundingCalcTime, mtp.AmmPoolId)

	if mtp.Position == types.Position_LONG {
		takeAmountCustodyAmount := types.CalcTakeAmount(mtp.Custody, longRate)
		if !takeAmountCustodyAmount.IsPositive() {
			return nil
		}

		// increase fees collected
		err := pool.UpdateFeesCollected(mtp.CustodyAsset, takeAmountCustodyAmount, true)
		if err != nil {
			return err
		}

		// update mtp custody
		mtp.Custody = mtp.Custody.Sub(takeAmountCustodyAmount)

		// add payment to total funding fee paid in custody asset
		mtp.FundingFeePaidCustody = mtp.FundingFeePaidCustody.Add(takeAmountCustodyAmount)

		// update pool custody balance
		err = pool.UpdateCustody(mtp.CustodyAsset, takeAmountCustodyAmount, false, mtp.Position)
		if err != nil {
			return err
		}
	} else {
		takeAmountLiabilityAmount := types.CalcTakeAmount(mtp.Liabilities, shortRate)
		if !takeAmountLiabilityAmount.IsPositive() {
			return nil
		}

		// increase fees collected
		// Note: fees is collected in liabilities asset
		err := pool.UpdateFeesCollected(mtp.LiabilitiesAsset, takeAmountLiabilityAmount, true)
		if err != nil {
			return err
		}

		tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
		if err != nil {
			return err
		}

		// should be done in custody
		// short -> usdc
		// long -> custody
		// For short, takeAmountLiabilityAmount is in trading asset, need to convert to custody asset which is in usdc
		custodyAmt := takeAmountLiabilityAmount.ToLegacyDec().Mul(tradingAssetPrice).TruncateInt()

		// update mtp custody
		mtp.Custody = mtp.Custody.Sub(custodyAmt)
		// add payment to total funding fee paid in custody asset
		mtp.FundingFeePaidCustody = mtp.FundingFeePaidCustody.Add(custodyAmt)

		// update pool custody balance
		err = pool.UpdateCustody(mtp.CustodyAsset, custodyAmt, false, mtp.Position)
		if err != nil {
			return err
		}
	}

	return nil
}
