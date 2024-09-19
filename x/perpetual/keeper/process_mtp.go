package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckAndLiquidateUnhealthyPosition(ctx sdk.Context, mtp *types.MTP, pool types.Pool, ammPool ammtypes.Pool, baseCurrency string, baseCurrencyDecimal uint64) error {
	var err error
	// Handle toPay
	err = k.HandleToPay(ctx)
	if err != nil {
		ctx.Logger().Error(err.Error())
	}

	// update mtp take profit liabilities
	// calculate mtp take profit liabilities, delta x_tp_l = delta y_tp_c * current price (take profit liabilities = take profit custody * current price)
	mtp.TakeProfitLiabilities, err = k.CalcMTPTakeProfitLiability(ctx, mtp, baseCurrency)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error calculating mtp take profit liabilities: %s", mtp.String()))
	}
	// calculate and update take profit borrow rate
	mtp.TakeProfitBorrowRate, err = k.CalcMTPTakeProfitBorrowRate(ctx, mtp)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error calculating mtp take profit borrow rate: %s", mtp.String()))
	}
	// Handle Borrow Interest if within epoch position
	if _, err := k.SettleBorrowInterest(ctx, mtp, &pool, ammPool); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error handling borrow interest payment: %s", mtp.CollateralAsset))
	}
	h, err := k.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error updating mtp health: %s", mtp.String()))
	}
	mtp.MtpHealth = h

	toPay, err := k.SettleFunding(ctx, mtp, &pool, ammPool, baseCurrency)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error handling funding fee: %s", mtp.CollateralAsset))
	}

	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	// check MTP health against threshold
	safetyFactor := k.GetSafetyFactor(ctx)
	var mustForceClose bool = false

	if mtp.MtpHealth.LTE(safetyFactor) {
		// flag position as must force close
		mustForceClose = true
	} else {
		ctx.Logger().Debug(errors.Wrap(types.ErrMTPHealthy, "skipping executing force close because mtp is healthy").Error())
	}

	entry, found := k.assetProfileKeeper.GetEntryByDenom(ctx, mtp.CustodyAsset)
	if !found {
		ctx.Logger().Error(errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", mtp.CustodyAsset).Error())
	}
	custodyAssetDecimals := entry.Decimals

	oneToken := math.NewIntFromBigInt(math.LegacyNewDec(10).Power(uint64(custodyAssetDecimals)).TruncateInt().BigInt())

	assetPrice, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.CustodyAsset, oneToken), baseCurrency, ammPool)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error estimating swap: %s", mtp.CustodyAsset))
	}

	// divide assetPrice by 10^baseCurrencyDecimal to get the actual price in decimal
	assetPriceDec := math.LegacyNewDecFromBigInt(assetPrice.BigInt()).Quo(math.LegacyNewDec(10).Power(uint64(baseCurrencyDecimal)))

	if types.ReachedTakeProfitPrice(mtp, assetPriceDec) {
		// flag position as must force close
		mustForceClose = true
	} else {
		ctx.Logger().Debug(fmt.Sprintf("skipping force close on position %s because take profit price %s <> %s", mtp.String(), mtp.TakeProfitPrice.String(), assetPrice.String()))
	}

	// if flag is false, then skip force close
	if !mustForceClose {
		return nil
	}

	var repayAmount math.Int
	switch mtp.Position {
	case types.Position_LONG:
		repayAmount, err = k.ForceCloseLong(ctx, mtp, &pool, true, baseCurrency)
	case types.Position_SHORT:
		repayAmount, err = k.ForceCloseShort(ctx, mtp, &pool, true, baseCurrency)
	default:
		return errors.Wrap(types.ErrInvalidPosition, fmt.Sprintf("invalid position type: %s", mtp.Position))
	}

	if err == nil {
		// Emit event if position was closed
		k.EmitForceClose(ctx, mtp, repayAmount, "")
	} else {
		return errors.Wrap(err, "error executing force close")
	}

	senderAddress, _ := sdk.AccAddressFromBech32(mtp.Address)
	found = k.DoesMTPExist(ctx, senderAddress, mtp.Id)
	empty := sdk.Coin{}
	if !found && toPay != empty {
		k.SetToPay(ctx, &types.ToPay{
			AssetDenom:   toPay.Denom,
			AssetBalance: toPay.Amount,
			Address:      senderAddress.String(),
			Id:           mtp.Id,
		})
	}

	return nil
}

func (k Keeper) CheckAndCloseAtStopLoss(ctx sdk.Context, mtp *types.MTP, pool types.Pool, ammPool ammtypes.Pool, baseCurrency string, baseCurrencyDecimal uint64) error {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()

	tradingAssetPrice, found := k.oracleKeeper.GetAssetPrice(ctx, mtp.TradingAsset)
	if !found {
		return fmt.Errorf("asset price not found")
	}

	if mtp.Position == types.Position_LONG {
		underStopLossPrice := !mtp.StopLossPrice.IsNil() && tradingAssetPrice.Price.LTE(mtp.StopLossPrice)
		if !underStopLossPrice {
			return fmt.Errorf("mtp stop loss price is not <=  token price")
		}
	} else {
		underStopLossPrice := !mtp.StopLossPrice.IsNil() && tradingAssetPrice.Price.GTE(mtp.StopLossPrice)
		if !underStopLossPrice {
			return fmt.Errorf("mtp stop loss price is not =>  token price")
		}
	}

	var (
		repayAmount math.Int
		err         error
	)
	switch mtp.Position {
	case types.Position_LONG:
		repayAmount, err = k.ForceCloseLong(ctx, mtp, &pool, true, baseCurrency)
	case types.Position_SHORT:
		repayAmount, err = k.ForceCloseShort(ctx, mtp, &pool, true, baseCurrency)
	default:
		return errors.Wrap(types.ErrInvalidPosition, fmt.Sprintf("invalid position type: %s", mtp.Position))
	}

	if err == nil {
		// Emit event if position was closed
		k.EmitForceClose(ctx, mtp, repayAmount, "")
	} else {
		return errors.Wrap(err, "error executing force close")
	}

	return nil
}

func (k Keeper) HandleToPay(ctx sdk.Context) error {
	toPays := k.GetAllToPayStore(ctx)

	if len(toPays) == 0 {
		return nil
	}
	// get funding fee collection address
	fundingFeeCollectionAddress := k.GetFundingFeeCollectionAddress(ctx)

	for _, toPay := range toPays {
		balance := k.bankKeeper.GetBalance(ctx, fundingFeeCollectionAddress, toPay.AssetDenom)
		if balance.Amount.LT(toPay.AssetBalance) {
			break
		} else {
			// transfer funding fee amount to mtp address
			if err := k.bankKeeper.SendCoins(ctx, fundingFeeCollectionAddress, sdk.MustAccAddressFromBech32(toPay.Address), sdk.NewCoins(sdk.NewCoin(toPay.AssetDenom, toPay.AssetBalance))); err != nil {
				return err
			}
			err := k.DeleteToPay(ctx, sdk.MustAccAddressFromBech32(toPay.Address), toPay.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
