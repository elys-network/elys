package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CalcReturnAmount(ctx sdk.Context, mtp types.MTP, pool types.Pool, ammPool ammtypes.Pool, repayAmount math.Int, amount math.Int, baseCurrency string) (returnAmount math.Int, err error) {
	Liabilities := mtp.Liabilities
	BorrowInterestUnpaid := mtp.BorrowInterestUnpaidCollateral

	if mtp.BorrowInterestUnpaidCollateral.IsPositive() {
		if mtp.Position == types.Position_SHORT {
			// swap to trading asset
			unpaidCollateralIn := sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral)
			C, err := k.CloseEstimationChecker.EstimateSwapGivenOut(ctx, unpaidCollateralIn, mtp.TradingAsset, ammPool)
			if err != nil {
				return sdk.ZeroInt(), err
			}

			BorrowInterestUnpaid = C
		} else if mtp.CollateralAsset != baseCurrency {
			// swap to base currency
			unpaidCollateralIn := sdk.NewCoin(mtp.CollateralAsset, mtp.BorrowInterestUnpaidCollateral)
			C, err := k.CloseEstimationChecker.EstimateSwapGivenOut(ctx, unpaidCollateralIn, baseCurrency, ammPool)
			if err != nil {
				return sdk.ZeroInt(), err
			}

			BorrowInterestUnpaid = C
		}
	}

	// Reminder:
	// if long both repay amount and liablities are collateral asset
	// if short both repay amount and liablities are trading asset

	have := repayAmount
	owe := Liabilities.Add(BorrowInterestUnpaid).Mul(amount.Quo(mtp.Custody))

	if have.LT(owe) {
		// v principle liability; x excess liability
		returnAmount = sdk.ZeroInt()
	} else {
		// can afford both
		returnAmount = have.Sub(owe)
	}

	// returnAmount is so far in base currency if long or trading asset if short, now should convert it to collateralAsset in order to return
	if returnAmount.IsPositive() {
		if mtp.Position == types.Position_SHORT {
			// swap to collateral asset
			amtTokenIn := sdk.NewCoin(mtp.TradingAsset, returnAmount)
			C, err := k.CloseEstimationChecker.EstimateSwapGivenOut(ctx, amtTokenIn, mtp.CollateralAsset, ammPool)
			if err != nil {
				return sdk.ZeroInt(), err
			}

			returnAmount = C
		} else if mtp.CollateralAsset != baseCurrency {
			// swap to collateral asset
			amtTokenIn := sdk.NewCoin(baseCurrency, returnAmount)
			C, err := k.CloseEstimationChecker.EstimateSwapGivenOut(ctx, amtTokenIn, mtp.CollateralAsset, ammPool)
			if err != nil {
				return sdk.ZeroInt(), err
			}

			returnAmount = C
		}
	}

	return returnAmount, nil
}
