package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) ProcessOpen(ctx sdk.Context, mtp *types.MTP, proxyLeverage osmomath.BigDec, collateralAmountDec osmomath.BigDec, poolId uint64, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	// Fetch the pool associated with the given pool ID.
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolDoesNotExist, "pool id %d", poolId)
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "amm pool id %d", poolId)
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := collateralAmountDec.Mul(proxyLeverage).Dec().TruncateInt()

	// Calculate custody amount
	// LONG: if collateral asset is trading asset then custodyAmount = leveragedAmount else if it collateral asset is usdc, we swap it to trading asset below
	// SHORT: collateralAsset is always usdc, and custody has to be in base currency, so custodyAmount = leveragedAmount
	custodyAmount := leveragedAmount

	switch msg.Position {
	case types.Position_LONG:
		// If collateral is not base currency, calculate the borrowing amount in base currency and check the balance
		if mtp.CollateralAsset != baseCurrency {
			custodyAmtToken := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
			borrowingAmount, _, _, err := k.EstimateSwapGivenOut(ctx, custodyAmtToken, baseCurrency, ammPool, mtp.Address)
			if err != nil {
				return nil, err
			}
			if !types.HasSufficientPoolBalance(ammPool, baseCurrency, borrowingAmount) {
				return nil, errorsmod.Wrap(types.ErrBorrowTooHigh, borrowingAmount.String())
			}
		} else {
			if !types.HasSufficientPoolBalance(ammPool, mtp.CollateralAsset, leveragedAmount) {
				return nil, errorsmod.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
			}
		}

		// If position is long, calculate custody amount in custody asset
		if mtp.CollateralAsset == baseCurrency {
			leveragedAmtTokenIn := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
			custodyAmount, _, _, err = k.EstimateSwapGivenIn(ctx, leveragedAmtTokenIn, mtp.CustodyAsset, ammPool, mtp.Address)
			if err != nil {
				return nil, err
			}
		}
	case types.Position_SHORT:
		if mtp.CollateralAsset != baseCurrency {
			return nil, errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "collateral must be base currency")
		}

		// check the balance
		if !types.HasSufficientPoolBalance(ammPool, mtp.CustodyAsset, custodyAmount) {
			return nil, errorsmod.Wrap(types.ErrBorrowTooHigh, custodyAmount.String())
		}
	default:
		return nil, errorsmod.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	// Ensure the AMM pool has enough balance.
	if !types.HasSufficientPoolBalance(ammPool, mtp.CustodyAsset, custodyAmount) {
		return nil, errorsmod.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
	}

	// Borrow the asset the user wants to long.
	err = k.Borrow(ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, proxyLeverage, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	if err = k.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Update the MTP health.
	mtpHealth, err := k.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}
	mtp.MtpHealth = mtpHealth.Dec()

	// Check if the MTP is unhealthy
	safetyFactor := k.GetSafetyFactor(ctx)
	if mtp.MtpHealth.LTE(safetyFactor) {
		return nil, errorsmod.Wrapf(types.ErrMTPUnhealthy, "(MtpHealth: %s)", mtp.MtpHealth.String())
	}

	// Set stop loss price
	// If consolidating or adding collateral, this needs to be calculated again
	stopLossPrice := msg.StopLossPrice
	if msg.StopLossPrice.IsNil() || msg.StopLossPrice.IsZero() {
		stopLossPrice = k.GetLiquidationPrice(ctx, *mtp).Dec()
	}
	mtp.StopLossPrice = stopLossPrice

	// Set MTP
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return nil, err
	}

	return mtp, nil
}
