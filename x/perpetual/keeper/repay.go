package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

// Repay ammPool has to be pointer because RemoveFromPoolBalance updates pool assets
func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool, returnAmount math.Int, payingLiabilities math.Int, closingRatio math.LegacyDec, perpFees *types.PerpetualFees, repayAmount math.Int, isLiquidation bool) error {
	if returnAmount.IsPositive() {
		ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
		if err != nil {
			return err
		}

		// send fees to masterchef and taker collection address
		totalFees, err := k.SendFeesToPoolRevenueAndTakerCollection(ctx, ammPoolAddr, mtp.Address, repayAmount, mtp.CustodyAsset, ammPool, perpFees, returnAmount)
		if err != nil {
			return err
		}

		// to prevent zero return amount
		if totalFees.LT(returnAmount) {
			returnCoins := sdk.NewCoins(sdk.NewCoin(mtp.CustodyAsset, returnAmount.Sub(totalFees)))
			returnReceiver := mtp.GetAccountAddress()
			if isLiquidation && mtp.PartialLiquidationDone {
				returnReceiver = sdk.MustAccAddressFromBech32(k.tierKeeper.GetMasterChefParams(ctx).ProtocolRevenueAddress)
			}

			err = k.SendFromAmmPool(ctx, ammPool, returnReceiver, returnCoins)
			if err != nil {
				return err
			}
		}
	}

	if isLiquidation && !mtp.PartialLiquidationDone {
		mtp.PartialLiquidationDone = true
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
