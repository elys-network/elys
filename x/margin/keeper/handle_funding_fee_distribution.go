package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

// HandleFundingFeeDistribution handles funding fee distribution
func (k Keeper) HandleFundingFeeDistribution(ctx sdk.Context, mtps []*types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) error {
	// get funding rate
	fundingRate := pool.FundingRate

	// if funding rate is zero, return
	if fundingRate.IsZero() {
		return nil
	}

	// account liabilities from long position
	liabilitiesLong := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsLong {
		liabilitiesLong = liabilitiesLong.Add(asset.Liabilities)
	}

	// account liabilities from short position
	liabilitiesShort := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsShort {
		liabilitiesShort = liabilitiesShort.Add(asset.Liabilities)
	}

	// get funding fee collection address
	fundingFeeCollectionAddress := k.GetFundingFeeCollectionAddress(ctx)

	// get base currency balance
	balance := k.bankKeeper.GetBalance(ctx, fundingFeeCollectionAddress, baseCurrency)

	// if balance is zero, return
	if balance.IsZero() {
		return nil
	}

	for _, mtp := range mtps {
		// if funding rate is negative and mtp position is short or funding rate is positive and mtp position is long, return
		if (fundingRate.IsNegative() && mtp.Position == types.Position_SHORT) || (fundingRate.IsPositive() && mtp.Position == types.Position_LONG) {
			return nil
		}

		// get mtp address
		mtpAddress, err := sdk.AccAddressFromBech32(mtp.Address)
		if err != nil {
			return err
		}

		// calc funding fee share
		fundingFeeShare := sdk.ZeroDec()
		if fundingRate.IsNegative() && mtp.Position == types.Position_LONG {
			fundingFeeShare = sdk.NewDecFromInt(mtp.Liabilities).Quo(sdk.NewDecFromInt(liabilitiesLong))
		}
		if fundingRate.IsPositive() && mtp.Position == types.Position_SHORT {
			fundingFeeShare = sdk.NewDecFromInt(mtp.Liabilities).Quo(sdk.NewDecFromInt(liabilitiesShort))
		}

		// if funding fee share is zero, skip mtp
		if fundingFeeShare.IsZero() {
			continue
		}

		// calculate funding fee amount
		fundingFeeAmount := sdk.NewCoin(baseCurrency, sdk.NewDecFromInt(balance.Amount).Mul(fundingFeeShare).TruncateInt())

		// transfer funding fee amount to mtp address
		if err := k.bankKeeper.SendCoins(ctx, fundingFeeCollectionAddress, mtpAddress, sdk.NewCoins(fundingFeeAmount)); err != nil {
			return err
		}

		// update received funding fee accounting buckets
		for custodyIndex, _ := range mtp.Custodies {
			for collateralIndex, collateral := range mtp.Collaterals {
				// Swap the take amount to collateral asset
				fundingFeeCollateralAmount, err := k.EstimateSwap(ctx, fundingFeeAmount, collateral.Denom, ammPool)
				if err != nil {
					return err
				}

				// Create the take amount coin
				fundingFeeCollateral := sdk.NewCoin(collateral.Denom, fundingFeeCollateralAmount)

				// add payment to total funding fee paid in collateral asset
				mtp.FundingFeeReceivedCollaterals[collateralIndex] = mtp.FundingFeePaidCollaterals[collateralIndex].Add(fundingFeeCollateral)

				// add payment to total funding fee paid in custody asset
				mtp.FundingFeeReceivedCustodies[custodyIndex] = mtp.FundingFeePaidCustodies[custodyIndex].Add(fundingFeeAmount)
			}
		}
	}

	return nil
}
