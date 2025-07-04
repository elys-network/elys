package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

// Repay ammPool has to be pointer because RemoveFromPoolBalance updates pool assets
func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool, returnAmount math.Int, payingLiabilities math.Int, closingRatio math.LegacyDec, baseCurrency string) error {
	if returnAmount.IsPositive() {
		ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
		if err != nil {
			return err
		}
		// send fees to masterchef and taker collection address
		totalFees, err := k.SendFeesToMasterchefAndTakerCollection(ctx, ammPoolAddr, mtp.Address, returnAmount, mtp.CustodyAsset, ammPool)
		if err != nil {
			return err
		}

		returnCoins := sdk.NewCoins(sdk.NewCoin(mtp.CustodyAsset, returnAmount.Sub(totalFees)))
		err = k.SendFromAmmPool(ctx, ammPool, mtp.GetAccountAddress(), returnCoins)
		if err != nil {
			return err
		}
	}

	mtp.Liabilities = mtp.Liabilities.Sub(payingLiabilities)

	closingCustodyAmount := closingRatio.MulInt(mtp.Custody).TruncateInt()
	mtp.Custody = mtp.Custody.Sub(closingCustodyAmount)

	reducingCollateralAmt := closingRatio.MulInt(mtp.Collateral).TruncateInt()
	mtp.Collateral = mtp.Collateral.Sub(reducingCollateralAmt)

	err := pool.UpdateCustody(mtp.CustodyAsset, closingCustodyAmount, false, mtp.Position)
	if err != nil {
		return err
	}

	err = pool.UpdateLiabilities(mtp.LiabilitiesAsset, payingLiabilities, false, mtp.Position)
	if err != nil {
		return err
	}

	err = pool.UpdateCollateral(mtp.CollateralAsset, reducingCollateralAmt, false, mtp.Position)
	if err != nil {
		return err
	}

	// This is for accounting purposes, mtp.Custody gets reduced by borrowInterestPaymentCustody and funding fee. so msg.Amount is greater than mtp.Custody here. So if it's negative it should be closed
	if mtp.Custody.IsZero() || mtp.Custody.IsNegative() {
		k.DestroyMTP(ctx, *mtp)
	} else {
		// update mtp health
		mtp.MtpHealth, err = k.GetMTPHealth(ctx, *mtp)
		if err != nil {
			return err
		}
		err = k.SetMTP(ctx, mtp)
		if err != nil {
			return err
		}
	}

	k.SetPool(ctx, *pool)

	return nil
}
