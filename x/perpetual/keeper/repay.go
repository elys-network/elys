package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// Repay ammPool has to be pointer because RemoveFromPoolBalance updates pool assets
func (k Keeper) Repay(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool, returnAmount math.Int, payingLiabilities math.Int, closingRatio math.LegacyDec, perpFees *types.PerpetualFees, repayAmount math.Int, isLiquidation bool, tradingAssetDenomPrice osmomath.BigDec) (sdk.Coin, error) {
	var collateralToAdd sdk.Coin
	if returnAmount.IsPositive() {
		ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
		if err != nil {
			return sdk.Coin{}, err
		}

		// send fees to masterchef and taker collection address
		totalFees, err := k.SendFeesToPoolRevenueAndTakerCollection(ctx, ammPoolAddr, mtp.Address, repayAmount, mtp.CustodyAsset, ammPool, perpFees, returnAmount)
		if err != nil {
			return sdk.Coin{}, err
		}

		// to prevent zero return amount
		if totalFees.LT(returnAmount) {
			returnCustodyCoins := sdk.NewCoins(sdk.NewCoin(mtp.CustodyAsset, returnAmount.Sub(totalFees)))

			returnCoins := returnCustodyCoins

			returnReceiver := mtp.GetAccountAddress()

			if isLiquidation {
				if mtp.PartialLiquidationDone {
					returnReceiver = sdk.MustAccAddressFromBech32(k.tierKeeper.GetMasterChefParams(ctx).ProtocolRevenueAddress)
				} else {
					// This means this is first liquidation
					collateralToAdd = returnCoins[0]
					// If this is the first liquidation due to poor health, then we return the amount to the user then later we use that amount as collateral
					// so we need to convert return amount to collateral asset denom.
					// For AMM, in case 2 and 3 below, there won't be change in balance
					// in case 1, asset gets swapped but TVL will be same

					// 1. Long & collateral is USDC, liabilities is base currency: Convert to return amount to USDC with oracle price
					// 2. Long & collateral is a trading asset thus same as custody asset, no need to do anything
					// 3. Short: Custody is already USDC - same as collateral
					// For long custody asset is trading asset
					if mtp.IsLong() && mtp.CollateralAsset == mtp.LiabilitiesAsset {
						amount := tradingAssetDenomPrice.Mul(osmomath.BigDecFromSDKInt(returnCustodyCoins[0].Amount)).Dec().TruncateInt()
						collateralToAdd = sdk.NewCoin(mtp.CollateralAsset, amount)
						returnCoins = sdk.NewCoins(collateralToAdd)
					}
				}
			}

			err = k.SendFromAmmPool(ctx, ammPool, returnReceiver, returnCoins)
			if err != nil {
				return sdk.Coin{}, err
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
		return sdk.Coin{}, err
	}

	err = pool.UpdateLiabilities(mtp.LiabilitiesAsset, payingLiabilities, false, mtp.Position)
	if err != nil {
		return sdk.Coin{}, err
	}

	err = pool.UpdateCollateral(mtp.CollateralAsset, reducingCollateralAmt, false, mtp.Position)
	if err != nil {
		return sdk.Coin{}, err
	}

	// This is for accounting purposes, mtp.Custody gets reduced by borrowInterestPaymentCustody and funding fee. so msg.Amount is greater than mtp.Custody here. So if it's negative it should be closed
	if mtp.Custody.IsZero() || mtp.Custody.IsNegative() {
		k.DestroyMTP(ctx, *mtp)
	} else {
		// update mtp health
		mtp.MtpHealth, err = k.GetMTPHealth(ctx, *mtp)
		if err != nil {
			return sdk.Coin{}, err
		}
		err = k.SetMTP(ctx, mtp)
		if err != nil {
			return sdk.Coin{}, err
		}
	}

	k.SetPool(ctx, *pool)

	return collateralToAdd, nil
}
