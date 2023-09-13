package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, repayAmount sdk.Int, takeFundPayment bool, collateralAsset string) error {
	collateralIndex, _ := k.GetMTPAssetIndex(mtp, collateralAsset, "")
	// nolint:staticcheck,ineffassign
	returnAmount, debtP, debtI := sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
	Liabilities := mtp.Liabilities
	InterestUnpaidCollateral := mtp.InterestUnpaidCollaterals[collateralIndex]

	if collateralAsset != ptypes.USDC {
		// swap to usdc
		unpaidCollateralIn := sdk.NewCoin(mtp.CollateralAssets[collateralIndex], mtp.InterestUnpaidCollaterals[collateralIndex])
		C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, ptypes.USDC, ammPool)
		if err != nil {
			return err
		}

		InterestUnpaidCollateral = C
	}

	var err error
	mtp.MtpHealth, err = k.UpdateMTPHealth(ctx, *mtp, ammPool)
	if err != nil {
		return err
	}

	have := repayAmount
	owe := Liabilities.Add(InterestUnpaidCollateral)

	if have.LT(Liabilities) {
		//can't afford principle liability
		returnAmount = sdk.ZeroInt()
		debtP = Liabilities.Sub(have)
		debtI = InterestUnpaidCollateral
	} else if have.LT(owe) {
		// v principle liability; x excess liability
		returnAmount = sdk.ZeroInt()
		debtP = sdk.ZeroInt()
		debtI = Liabilities.Add(InterestUnpaidCollateral).Sub(have)
	} else {
		// can afford both
		returnAmount = have.Sub(Liabilities).Sub(InterestUnpaidCollateral)
		debtP = sdk.ZeroInt()
		debtI = sdk.ZeroInt()
	}

	if !returnAmount.IsZero() {
		actualReturnAmount := returnAmount
		if takeFundPayment {
			takePercentage := k.GetForceCloseFundPercentage(ctx)

			fundAddr := k.GetForceCloseFundAddress(ctx)
			takeAmount, err := k.TakeFundPayment(ctx, returnAmount, ptypes.USDC, takePercentage, fundAddr, &ammPool)
			if err != nil {
				return err
			}
			actualReturnAmount = returnAmount.Sub(takeAmount)
			if !takeAmount.IsZero() {
				k.EmitFundPayment(ctx, mtp, takeAmount, ptypes.USDC, types.EventRepayFund)
			}
		}

		// actualReturnAmount is so far in usdc, now should convert it to collateralAsset in order to return
		if !actualReturnAmount.IsZero() {
			if collateralAsset != ptypes.USDC {
				// swap to usdc
				amtTokenIn := sdk.NewCoin(ptypes.USDC, actualReturnAmount)
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
	// because so far returnAmount is in usdc.
	if collateralAsset != ptypes.USDC {
		// swap to usdc
		amtTokenIn := sdk.NewCoin(ptypes.USDC, returnAmount)
		C, err := k.EstimateSwapGivenOut(ctx, amtTokenIn, collateralAsset, ammPool)
		if err != nil {
			return err
		}

		returnAmount = C
	}

	err = pool.UpdateBalance(ctx, mtp.CollateralAssets[collateralIndex], returnAmount, false)
	if err != nil {
		return err
	}

	// Need to be checked by Caner for short
	// long position
	err = pool.UpdateLiabilities(ctx, ptypes.USDC, mtp.Liabilities, false)
	if err != nil {
		return err
	}

	// long position
	err = pool.UpdateUnsettledLiabilities(ctx, ptypes.USDC, debtI, true)
	if err != nil {
		return err
	}

	// long position
	err = pool.UpdateUnsettledLiabilities(ctx, ptypes.USDC, debtP, true)
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
