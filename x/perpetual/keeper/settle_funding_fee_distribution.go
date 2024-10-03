package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) SettleFundingFeeDistribution(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, baseCurrency string) (sdk.Coin, error) {
	// get mtp address
	mtpAddress, err := sdk.AccAddressFromBech32(mtp.Address)
	if err != nil {
		return sdk.Coin{}, err
	}

	uusdc, found := k.assetProfileKeeper.GetEntry(ctx, "uusdc")
	if !found {
		return sdk.Coin{}, nil
	}

	// Calculate liabilities for long and short assets using the separate helper function
	liabilitiesLong, err := k.CalcTotalLiabilities(ctx, pool.PoolAssetsLong, pool.AmmPoolId, uusdc.Denom)
	if err != nil {
		return sdk.Coin{}, nil
	}

	liabilitiesShort, err := k.CalcTotalLiabilities(ctx, pool.PoolAssetsShort, pool.AmmPoolId, uusdc.Denom)
	if err != nil {
		return sdk.Coin{}, nil
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
	long, short := k.GetFundingDistributionValue(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId)
	var totalFund sdkmath.LegacyDec
	// calc funding fee share
	var fundingFeeShare sdkmath.LegacyDec
	if mtp.Position == types.Position_LONG {
		// Ensure liabilitiesLong is not zero to avoid division by zero
		if liabilitiesLong.IsZero() {
			return sdk.Coin{}, types.ErrAmountTooLow
		}
		fundingFeeShare = sdkmath.LegacyNewDecFromInt(mtp.Liabilities).Quo(sdkmath.LegacyNewDecFromInt(liabilitiesLong))
		totalFund = short
	} else {
		// Ensure liabilitiesShort is not zero to avoid division by zero
		if liabilitiesShort.IsZero() {
			return sdk.Coin{}, types.ErrAmountTooLow
		}
		fundingFeeShare = sdkmath.LegacyNewDecFromInt(mtp.Liabilities).Quo(sdkmath.LegacyNewDecFromInt(liabilitiesShort))
		totalFund = long
	}

	// if funding fee share is zero, skip mtp
	if fundingFeeShare.IsZero() {
		return sdk.Coin{}, nil
	}

	// calculate funding fee amount
	fundingFeeAmount := sdk.NewCoin(baseCurrency, totalFund.Mul(fundingFeeShare).TruncateInt())
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

	// add payment to total funding fee paid in collateral asset
	mtp.FundingFeeReceivedCollateral = mtp.FundingFeeReceivedCollateral.Add(fundingFeeCollateralAmount)
	// add payment to total funding fee paid in custody asset
	mtp.FundingFeeReceivedCustody = mtp.FundingFeeReceivedCustody.Add(fundingFeeAmount.Amount)

	return toPay, nil
}
