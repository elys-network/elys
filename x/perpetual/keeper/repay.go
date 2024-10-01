package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, repayAmount math.Int, takeFundPayment bool, amount math.Int, baseCurrency string) error {
	// nolint:staticcheck,ineffassign
	returnAmount := sdk.ZeroInt()
	Liabilities := mtp.Liabilities
	BorrowInterestUnpaid := mtp.BorrowInterestUnpaidCollateral

	if mtp.BorrowInterestUnpaidCollateral.IsPositive() {
		if mtp.Position == types.Position_SHORT {
			// swap to trading asset
			unpaidCollateralIn := sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral)
			C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, mtp.TradingAsset, ammPool)
			if err != nil {
				return err
			}

			BorrowInterestUnpaid = C
		} else if mtp.CollateralAsset != baseCurrency {
			// swap to base currency
			unpaidCollateralIn := sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral)
			C, err := k.EstimateSwapGivenOut(ctx, unpaidCollateralIn, baseCurrency, ammPool)
			if err != nil {
				return err
			}

			BorrowInterestUnpaid = C
		}
	}

	var err error
	mtp.MtpHealth, err = k.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return err
	}

	// Reminder:
	// if long both repay amount and liablities are collateral asset
	// if short both repay amount and liablities are trading asset

	have := repayAmount
	owe := Liabilities.Add(BorrowInterestUnpaid).Mul(amount).Quo(mtp.Custody)

	if have.LT(Liabilities) {
		// can't afford principle liability
		returnAmount = sdk.ZeroInt()
	} else if have.LT(owe) {
		// v principle liability; x excess liability
		returnAmount = sdk.ZeroInt()
	} else {
		// can afford both
		returnAmount = have.Sub(Liabilities).Sub(BorrowInterestUnpaid)
		mtp.Liabilities = mtp.Liabilities.Sub(Liabilities.Mul(amount).Quo(mtp.Custody))
	}

	if returnAmount.IsPositive() {
		actualReturnAmount := returnAmount
		if takeFundPayment {
			takePercentage := k.GetForceCloseFundPercentage(ctx)

			fundAddr := k.GetForceCloseFundAddress(ctx)
			takeAmount := sdk.ZeroInt()
			if mtp.Position == types.Position_LONG {
				takeAmount, err = k.TakeFundPayment(ctx, returnAmount, baseCurrency, takePercentage, fundAddr, &ammPool)
				if err != nil {
					return err
				}
			} else if mtp.Position == types.Position_SHORT {
				takeAmount, err = k.TakeFundPayment(ctx, returnAmount, mtp.TradingAsset, takePercentage, fundAddr, &ammPool)
				if err != nil {
					return err
				}
			}
			actualReturnAmount = returnAmount.Sub(takeAmount)
			if !takeAmount.IsZero() {
				k.EmitFundPayment(ctx, mtp, takeAmount, baseCurrency, types.EventRepayFund)
			}
		}

		// actualReturnAmount is so far in base currency if long or trading asset if short, now should convert it to collateralAsset in order to return
		if actualReturnAmount.IsPositive() {
			if mtp.Position == types.Position_SHORT {
				// swap to collateral asset
				amtTokenIn := sdk.NewCoin(mtp.TradingAsset, actualReturnAmount)
				C, err := k.EstimateSwapGivenOut(ctx, amtTokenIn, mtp.CollateralAsset, ammPool)
				if err != nil {
					return err
				}

				actualReturnAmount = C
			} else if mtp.CollateralAsset != baseCurrency {
				// swap to collateral asset
				amtTokenIn := sdk.NewCoin(baseCurrency, actualReturnAmount)
				C, err := k.EstimateSwapGivenOut(ctx, amtTokenIn, mtp.CollateralAsset, ammPool)
				if err != nil {
					return err
				}

				actualReturnAmount = C
			}

			returnCoin := sdk.NewCoin(mtp.CollateralAsset, sdk.NewIntFromBigInt(actualReturnAmount.BigInt()))
			returnCoins := sdk.NewCoins(returnCoin)
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
	if mtp.Position == types.Position_SHORT {
		// swap to collateral asset
		amtTokenIn := sdk.NewCoin(mtp.TradingAsset, returnAmount)
		C, err := k.EstimateSwapGivenOut(ctx, amtTokenIn, mtp.CollateralAsset, ammPool)
		if err != nil {
			return err
		}

		returnAmount = C
	} else if mtp.CollateralAsset != baseCurrency {
		// swap to collateral asset
		amtTokenIn := sdk.NewCoin(baseCurrency, returnAmount)
		C, err := k.EstimateSwapGivenOut(ctx, amtTokenIn, mtp.CollateralAsset, ammPool)
		if err != nil {
			return err
		}

		returnAmount = C
	}

	err = pool.UpdateBalance(ctx, mtp.CollateralAsset, returnAmount, false, mtp.Position)
	if err != nil {
		return err
	}

	// long position
	err = pool.UpdateLiabilities(ctx, mtp.LiabilitiesAsset, mtp.Liabilities, false, mtp.Position)
	if err != nil {
		return err
	}

	err = pool.UpdateTakeProfitLiabilities(ctx, mtp.LiabilitiesAsset, mtp.TakeProfitLiabilities, false, mtp.Position)
	if err != nil {
		return err
	}

	err = pool.UpdateTakeProfitCustody(ctx, mtp.CustodyAsset, mtp.TakeProfitCustody, false, mtp.Position)
	if err != nil {
		return err
	}

	mtp.Custody = mtp.Custody.Sub(amount)
	// This is for accounting purposes, mtp.Custody gets reduced by borrowInterestPaymentCustody and funding fee. so msg.Amount is greater than mtp.Custody here. So if it's negative it should be closed
	if mtp.Custody.IsZero() || mtp.Custody.IsNegative() {
		err = k.DestroyMTP(ctx, mtp.GetAccountAddress(), mtp.Id)
		if err != nil {
			return err
		}
	} else {
		err = k.SetMTP(ctx, mtp)
		if err != nil {
			return err
		}
	}

	k.SetPool(ctx, *pool)

	return nil
}
