package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) SettleBorrowInterest(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) (sdkmath.Int, error) {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdkmath.ZeroInt(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	borrowInterestPaymentInt := k.GetBorrowInterest(ctx, mtp, ammPool)
	// pay interest+unpaid collateral amount
	finalBorrowInterestPayment, err := k.IncrementalBorrowInterestPayment(ctx, borrowInterestPaymentInt, mtp, pool, ammPool, baseCurrency)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "error executing incremental borrow interest payment").Error())
	}

	mtp.LastInterestCalcBlock = uint64(ctx.BlockHeight())
	mtp.LastInterestCalcTime = uint64(ctx.BlockTime().Unix())
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return finalBorrowInterestPayment, err
	}

	_, err = k.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	return finalBorrowInterestPayment, err
}

func (k Keeper) GetBorrowInterest(ctx sdk.Context, mtp *types.MTP, ammPool ammtypes.Pool) sdkmath.Int {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdkmath.ZeroInt()
	}
	baseCurrency := entry.Denom
	// Unpaid collateral
	unpaidCollateral := sdkmath.ZeroInt()
	if mtp.BorrowInterestUnpaidCollateral.IsPositive() {
		if mtp.CollateralAsset == baseCurrency {
			unpaidCollateral = unpaidCollateral.Add(mtp.BorrowInterestUnpaidCollateral)
		} else {
			// Liability is in base currency, so convert it to base currency
			unpaidCollateralIn := sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral)
			C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, baseCurrency, ammPool)
			if err != nil {
				return sdkmath.ZeroInt()
			}

			unpaidCollateral = unpaidCollateral.Add(C)
		}
	}
	sum := mtp.Liabilities.Add(unpaidCollateral)

	minBorrowInterestAmount := k.GetParams(ctx).MinBorrowInterestAmount
	// Get interest
	borrowInterestPayment := k.GetBorrowRate(ctx, mtp.LastInterestCalcBlock, mtp.AmmPoolId, sdkmath.LegacyNewDecFromInt(sum))
	return sdkmath.MaxInt(borrowInterestPayment.Mul(mtp.TakeProfitBorrowRate).TruncateInt(), minBorrowInterestAmount)
}
