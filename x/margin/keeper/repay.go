package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, repayAmount sdk.Int, takeFundPayment bool, collateralAsset string) error {
	collateralIndex, _ := k.GetMTPAssetIndex(mtp, collateralAsset, "")
	// nolint:staticcheck,ineffassign
	returnAmount := sdk.ZeroInt()
	Liabilities := mtp.Liabilities
	BorrowInterestUnpaidCollateral := mtp.BorrowInterestUnpaidCollaterals[collateralIndex]

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	if collateralAsset != baseCurrency {
		// swap to base currency
		unpaidCollateralIn := sdk.NewCoin(mtp.Collaterals[collateralIndex].Denom, mtp.BorrowInterestUnpaidCollaterals[collateralIndex].Amount)
		C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, baseCurrency, ammPool)
		if err != nil {
			return err
		}

		BorrowInterestUnpaidCollateral.Amount = C
	}

	var err error
	mtp.MtpHealth, err = k.UpdateMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return err
	}

	have := repayAmount
	owe := Liabilities.Add(BorrowInterestUnpaidCollateral.Amount)

	if have.LT(Liabilities) {
		//can't afford principle liability
		returnAmount = sdk.ZeroInt()
	} else if have.LT(owe) {
		// v principle liability; x excess liability
		returnAmount = sdk.ZeroInt()
	} else {
		// can afford both
		returnAmount = have.Sub(Liabilities).Sub(BorrowInterestUnpaidCollateral.Amount)
	}

	if !returnAmount.IsZero() {
		actualReturnAmount := returnAmount
		if takeFundPayment {
			takePercentage := k.GetForceCloseFundPercentage(ctx)

			fundAddr := k.GetForceCloseFundAddress(ctx)
			takeAmount, err := k.TakeFundPayment(ctx, returnAmount, baseCurrency, takePercentage, fundAddr, &ammPool)
			if err != nil {
				return err
			}
			actualReturnAmount = returnAmount.Sub(takeAmount)
			if !takeAmount.IsZero() {
				k.EmitFundPayment(ctx, mtp, takeAmount, baseCurrency, types.EventRepayFund)
			}
		}

		// actualReturnAmount is so far in base currency, now should convert it to collateralAsset in order to return
		if !actualReturnAmount.IsZero() {
			if collateralAsset != baseCurrency {
				// swap to base currency
				amtTokenIn := sdk.NewCoin(baseCurrency, actualReturnAmount)
				C, err := k.EstimateSwapGivenOut(ctx, amtTokenIn, collateralAsset, ammPool)
				if err != nil {
					return err
				}

				actualReturnAmount = C
			}

			var coins sdk.Coins
			returnCoin := sdk.NewCoin(collateralAsset, sdk.NewIntFromBigInt(actualReturnAmount.BigInt()))
			returnCoins := coins.Add(returnCoin)
			addr, err := sdk.AccAddressFromBech32(mtp.Address)
			if err != nil {
				return err
			}

			ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
			if err != nil {
				return err
			}

			err = k.bankKeeper.SendCoins(ctx, ammPoolAddr, addr, returnCoins)
			if err != nil {
				return err
			}
		}
	}

	// before updating collateral asset balance, we should convert returnAmount to collateralAsset
	// because so far returnAmount is in base currency.
	if collateralAsset != baseCurrency {
		// swap to base currency
		amtTokenIn := sdk.NewCoin(baseCurrency, returnAmount)
		C, err := k.EstimateSwapGivenOut(ctx, amtTokenIn, collateralAsset, ammPool)
		if err != nil {
			return err
		}

		returnAmount = C
	}

	err = pool.UpdateBalance(ctx, mtp.Collaterals[collateralIndex].Denom, returnAmount, false, mtp.Position)
	if err != nil {
		return err
	}

	// long position
	err = pool.UpdateLiabilities(ctx, baseCurrency, mtp.Liabilities, false, mtp.Position)
	if err != nil {
		return err
	}

	err = k.DestroyMTP(ctx, mtp.Address, mtp.Id)
	if err != nil {
		return err
	}

	k.SetPool(ctx, *pool)

	return nil
}
