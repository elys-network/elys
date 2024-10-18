package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) ProcessOpen(ctx sdk.Context, mtp *types.MTP, leverage sdk.Dec, eta sdk.Dec, collateralAmountDec sdk.Dec, poolId uint64, msg *types.MsgOpen, baseCurrency string, isBroker bool) (*types.MTP, error) {
	// Fetch the pool associated with the given pool ID.
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolDoesNotExist, "pool id %d", poolId)
	}

	// Check if the pool is enabled.
	if !pool.IsEnabled() {
		return nil, fmt.Errorf("disabled pool id %d", poolId)
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "amm pool id %d", poolId)
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := collateralAmountDec.Mul(leverage).TruncateInt()

	// Calculate custody amount
	// LONG: if collateral asset is trading asset then custodyAmount = leveragedAmount else if it collateral asset is usdc, we swap it to trading asset below
	// SHORT: collateralAsset is always usdc, and custody has to be in base currency, so custodyAmount = leveragedAmount
	custodyAmount := leveragedAmount

	switch msg.Position {
	case types.Position_LONG:
		// If collateral is not base currency, calculate the borrowing amount in base currency and check the balance
		if mtp.CollateralAsset != baseCurrency {
			custodyAmtToken := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
			borrowingAmount, _, err := k.EstimateSwapGivenOut(ctx, custodyAmtToken, baseCurrency, ammPool)
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
			custodyAmount, _, err = k.EstimateSwap(ctx, leveragedAmtTokenIn, mtp.CustodyAsset, ammPool)
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
	err = k.Borrow(ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, baseCurrency, isBroker)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	if err = k.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Take custody from the pool balance.
	if err = k.TakeInCustody(ctx, *mtp, &pool); err != nil {
		return nil, err
	}

	// Update the MTP health.
	lr, err := k.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// Set stop loss price
	mtp.StopLossPrice = msg.StopLossPrice

	// Set MTP
	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return nil, err
	}

	return mtp, nil
}
