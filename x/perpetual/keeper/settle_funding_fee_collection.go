package keeper

import (
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
	longRate, shortRate := k.GetFundingRate(ctx, mtp.LastFundingCalcBlock, mtp.LastFundingCalcTime, mtp.AmmPoolId)

	if mtp.Position == types.Position_LONG {
		takeAmountCustodyAmount := types.CalcTakeAmount(mtp.Custody, longRate)
		if !takeAmountCustodyAmount.IsPositive() {
			return nil
		}

		// increase fees collected
		err := pool.UpdateFeesCollected(ctx, mtp.CustodyAsset, takeAmountCustodyAmount, true)
		if err != nil {
			return err
		}

		// update mtp custody
		mtp.Custody = mtp.Custody.Sub(takeAmountCustodyAmount)

		// add payment to total funding fee paid in custody asset
		mtp.FundingFeePaidCustody = mtp.FundingFeePaidCustody.Add(takeAmountCustodyAmount)

		// update pool custody balance
		err = pool.UpdateCustody(ctx, mtp.CustodyAsset, takeAmountCustodyAmount, false, mtp.Position)
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
		err := pool.UpdateFeesCollected(ctx, mtp.LiabilitiesAsset, takeAmountLiabilityAmount, true)
		if err != nil {
			return err
		}

		// should be done in custody
		// short -> usdc
		// long -> custody
		custodyAmt, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.LiabilitiesAsset, takeAmountLiabilityAmount), mtp.CustodyAsset, ammPool)
		if err != nil {
			return err
		}

		// update mtp custody
		mtp.Custody = mtp.Custody.Sub(custodyAmt)
		// add payment to total funding fee paid in custody asset
		mtp.FundingFeePaidCustody = mtp.FundingFeePaidCustody.Add(custodyAmt)

		// update pool custody balance
		err = pool.UpdateCustody(ctx, mtp.CustodyAsset, custodyAmt, false, mtp.Position)
		if err != nil {
			return err
		}
	}

	// apply changes to mtp object
	err := k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	// apply changes to pool object
	k.SetPool(ctx, *pool)

	return nil
}
