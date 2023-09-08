package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool, repayAmount sdk.Int, takeFundPayment bool) error {
	// nolint:staticcheck,ineffassign
	returnAmount, debtP, debtI := sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
	Liabilities := mtp.Liabilities
	InterestUnpaidCollateral := mtp.InterestUnpaidCollateral

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
			takeAmount, err := k.TakeFundPayment(ctx, returnAmount, mtp.CollateralAsset, takePercentage, fundAddr, &ammPool)
			if err != nil {
				return err
			}
			actualReturnAmount = returnAmount.Sub(takeAmount)
			if !takeAmount.IsZero() {
				k.EmitFundPayment(ctx, mtp, takeAmount, mtp.CollateralAsset, types.EventRepayFund)
			}
		}

		if !actualReturnAmount.IsZero() {
			var coins sdk.Coins
			returnCoin := sdk.NewCoin(mtp.CollateralAsset, sdk.NewIntFromBigInt(actualReturnAmount.BigInt()))
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

	err = pool.UpdateBalance(ctx, mtp.CollateralAsset, returnAmount, false)
	if err != nil {
		return err
	}

	err = pool.UpdateLiabilities(ctx, mtp.CollateralAsset, mtp.Liabilities, false)
	if err != nil {
		return err
	}

	err = pool.UpdateUnsettledLiabilities(ctx, mtp.CollateralAsset, debtI, true)
	if err != nil {
		return err
	}

	err = pool.UpdateUnsettledLiabilities(ctx, mtp.CollateralAsset, debtP, true)
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
