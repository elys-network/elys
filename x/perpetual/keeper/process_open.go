package keeper

import (
	"fmt"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) ProcessOpen(ctx sdk.Context, pool *types.Pool, ammPool *ammtypes.Pool, mtp *types.MTP, proxyLeverage math.LegacyDec, poolId uint64, msg *types.MsgOpen, baseCurrency string) error {
	var err error
	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := proxyLeverage.MulInt(msg.Collateral.Amount).TruncateInt()

	// Calculate custody amount
	// LONG: if collateral asset is trading asset then custodyAmount = leveragedAmount else if it collateral asset is usdc, we swap it to trading asset below
	// SHORT: collateralAsset is always usdc, and custody has to be in base currency, so custodyAmount = leveragedAmount
	custodyAmount := leveragedAmount

	switch msg.Position {
	case types.Position_LONG:
		// If collateral is not base currency, calculate the borrowing amount in base currency and check the balance
		if mtp.CollateralAsset != baseCurrency {
			custodyAmtToken := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
			borrowingAmount, _, _, _, _, err := k.EstimateSwapGivenOut(ctx, custodyAmtToken, baseCurrency, *ammPool, mtp.Address)
			if err != nil {
				return err
			}
			if !types.HasSufficientPoolBalance(*ammPool, baseCurrency, borrowingAmount) {
				return errorsmod.Wrap(types.ErrBorrowTooHigh, borrowingAmount.String())
			}
		} else {
			if !types.HasSufficientPoolBalance(*ammPool, mtp.CollateralAsset, leveragedAmount) {
				return errorsmod.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
			}
		}

		// If position is long, calculate custody amount in custody asset
		if mtp.CollateralAsset == baseCurrency {
			leveragedAmtTokenIn := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
			custodyAmount, _, _, _, _, err = k.EstimateSwapGivenIn(ctx, leveragedAmtTokenIn, mtp.CustodyAsset, *ammPool, mtp.Address)
			if err != nil {
				return err
			}
		}
	case types.Position_SHORT:
		if mtp.CollateralAsset != baseCurrency {
			return errorsmod.Wrap(types.ErrInvalidCollateralAsset, "collateral must be base currency")
		}

		// check the balance
		if !types.HasSufficientPoolBalance(*ammPool, mtp.CustodyAsset, custodyAmount) {
			return errorsmod.Wrap(types.ErrBorrowTooHigh, custodyAmount.String())
		}
	default:
		return errorsmod.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	// Ensure the AMM pool has enough balance.
	if !types.HasSufficientPoolBalance(*ammPool, mtp.CustodyAsset, custodyAmount) {
		return errorsmod.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
	}

	// Borrow the asset the user wants to long.
	err = k.Borrow(ctx, msg.Collateral.Amount, custodyAmount, mtp, ammPool, pool, proxyLeverage, baseCurrency)
	if err != nil {
		return err
	}

	// Update the pool health.
	if err = k.UpdatePoolHealth(ctx, pool); err != nil {
		return err
	}

	// Update the MTP health.
	mtp.MtpHealth, err = k.GetMTPHealth(ctx, *mtp)
	if err != nil {
		return err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.GetSafetyFactor(ctx)
	if mtp.MtpHealth.LTE(safetyFactor) {
		return errorsmod.Wrapf(types.ErrMTPUnhealthy, "(MtpHealth: %s)", mtp.MtpHealth.String())
	}

	// Set stop loss price
	// If consolidating or adding collateral, this needs to be calculated again
	stopLossPrice := msg.StopLossPrice
	if msg.StopLossPrice.IsNil() || msg.StopLossPrice.IsZero() {
		stopLossPrice, err = k.GetLiquidationPrice(ctx, *mtp)
		if err != nil {
			return fmt.Errorf("failed to get liquidation price: %s", err.Error())
		}
	}
	mtp.StopLossPrice = stopLossPrice

	// calc and update open price
	err = k.GetAndSetOpenPrice(ctx, mtp, msg.Leverage.IsZero())
	if err != nil {
		return err
	}

	// Set MTP
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	return nil
}
