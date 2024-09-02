package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) SettleBorrowInterest(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) (math.Int, error) {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroInt(), errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	borrowInterestPaymentInt := k.GetBorrowInterest(ctx, mtp, pool, ammPool)
	// pay interest+unpaid collateral amount
	finalBorrowInterestPayment, err := k.IncrementalBorrowInterestPayment(ctx, borrowInterestPaymentInt, mtp, pool, ammPool, baseCurrency)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "error executing incremental borrow interest payment").Error())
	}

	_, err = k.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	return finalBorrowInterestPayment, err
}

func (k Keeper) GetBorrowInterest(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) math.Int {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.ZeroInt()
	}
	baseCurrency := entry.Denom
	// Unpaid collateral
	unpaidCollateral := sdk.ZeroInt()
	if mtp.CollateralAsset == baseCurrency {
		unpaidCollateral = unpaidCollateral.Add(mtp.BorrowInterestUnpaidCollateral)
	} else {
		// Liability is in base currency, so convert it to base currency
		unpaidCollateralIn := sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral)
		C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, baseCurrency, ammPool)
		if err != nil {
			return sdk.ZeroInt()
		}

		unpaidCollateral = unpaidCollateral.Add(C)
	}

	// Get interest
	borrowInterestPayment := k.GetBorrowRate(ctx, mtp.LastInterestCalcBlock, pool.AmmPoolId, mtp.LastInterestCalcTime, math.LegacyDec(mtp.Liabilities.Add(unpaidCollateral)))
	return borrowInterestPayment.Mul(mtp.TakeProfitBorrowRate).TruncateInt()
}
