package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// TODO: Think about funding rate algo, edge cases
func (k Keeper) SettleFundingFeeDistribution(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) (sdk.Coin, error) {
	// get funding rate
	fundingRate := k.GetFundingRate(ctx, mtp.LastFundingCalcBlock, mtp.AmmPoolId)

	// if funding rate is negative and mtp position is short or funding rate is positive and mtp position is long, return
	if (fundingRate.IsNegative() && mtp.Position == types.Position_SHORT) || (fundingRate.IsPositive() && mtp.Position == types.Position_LONG) {
		return sdk.Coin{}, nil
	}

	// get mtp address
	mtpAddress, err := sdk.AccAddressFromBech32(mtp.Address)
	if err != nil {
		return sdk.Coin{}, err
	}

	totalCustodyLong := sdk.ZeroInt()
	totalCustodyShort := sdk.ZeroInt()

	// account liabilities from long position
	liabilitiesLong := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsLong {
		liabilitiesLong = liabilitiesLong.Add(asset.Liabilities)
		totalCustodyLong = totalCustodyLong.Add(asset.Custody)
	}

	// account liabilities from short position
	liabilitiesShort := sdk.ZeroInt()
	for _, asset := range pool.PoolAssetsShort {
		liabilitiesShort = liabilitiesShort.Add(asset.Liabilities)
		totalCustodyShort = totalCustodyShort.Add(asset.Custody)
	}

	// get funding fee collection address
	fundingFeeCollectionAddress := k.GetFundingFeeCollectionAddress(ctx)

	// get base currency balance
	balance := k.bankKeeper.GetBalance(ctx, fundingFeeCollectionAddress, baseCurrency)

	// if balance is zero, return
	if balance.IsZero() {
		return sdk.Coin{}, nil
	}

	// Total fund collected should be
	var totalFund sdk.Int
	if fundingRate.IsNegative() {
		// short pays long
		totalFund = types.CalcTakeAmount(totalCustodyShort, fundingRate)
	} else {
		// long pays short
		totalFund = types.CalcTakeAmount(totalCustodyLong, fundingRate)
	}
	// calc funding fee share
	fundingFeeShare := sdk.ZeroDec()
	if fundingRate.IsNegative() && mtp.Position == types.Position_LONG {
		// Ensure liabilitiesLong is not zero to avoid division by zero
		if liabilitiesLong.IsZero() {
			return sdk.Coin{}, types.ErrAmountTooLow
		}
		fundingFeeShare = sdk.NewDecFromInt(mtp.Liabilities).Quo(sdk.NewDecFromInt(liabilitiesLong))
	}
	if fundingRate.IsPositive() && mtp.Position == types.Position_SHORT {
		// Ensure liabilitiesShort is not zero to avoid division by zero
		if liabilitiesShort.IsZero() {
			return sdk.Coin{}, types.ErrAmountTooLow
		}
		fundingFeeShare = sdk.NewDecFromInt(mtp.Liabilities).Quo(sdk.NewDecFromInt(liabilitiesShort))
	}

	// if funding fee share is zero, skip mtp
	if fundingFeeShare.IsZero() {
		return sdk.Coin{}, nil
	}

	// calculate funding fee amount
	fundingFeeAmount := sdk.NewCoin(baseCurrency, sdk.NewDecFromInt(totalFund).Mul(fundingFeeShare).TruncateInt())
	toPay := sdk.Coin{}

	if balance.Amount.LT(fundingFeeAmount.Amount) {
		toPay = fundingFeeAmount
	} else {
		// transfer funding fee amount to mtp address
		if err := k.bankKeeper.SendCoins(ctx, fundingFeeCollectionAddress, mtpAddress, sdk.NewCoins(fundingFeeAmount)); err != nil {
			return sdk.Coin{}, err
		}
	}

	// update received funding fee accounting buckets
	// Swap the take amount to collateral asset
	fundingFeeCollateralAmount, err := k.EstimateSwap(ctx, fundingFeeAmount, mtp.CollateralAsset, ammPool)
	if err != nil {
		return sdk.Coin{}, err
	}

	// TODO: What's the use of below fields ? should this be considered in MTP.Custody ?

	// add payment to total funding fee paid in collateral asset
	mtp.FundingFeeReceivedCollateral = mtp.FundingFeeReceivedCollateral.Add(fundingFeeCollateralAmount)

	// add payment to total funding fee paid in custody asset
	mtp.FundingFeeReceivedCustody = mtp.FundingFeeReceivedCustody.Add(fundingFeeAmount.Amount)

	mtp.LastFundingCalcBlock = uint64(ctx.BlockHeight())
	mtp.LastFundingCalcTime = uint64(ctx.BlockTime().Unix())

	return toPay, nil
}
