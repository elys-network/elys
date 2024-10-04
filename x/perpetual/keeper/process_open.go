package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) ProcessOpen(ctx sdk.Context, mtp *types.MTP, leverage sdk.Dec, eta sdk.Dec, collateralAmountDec sdk.Dec, poolId uint64, msg *types.MsgOpen, baseCurrency string, isBroker bool) (*types.MTP, error) {
	// Fetch the pool associated with the given pool ID.
	pool, found := k.OpenDefineAssetsChecker.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, mtp.TradingAsset)
	}

	// Check if the pool is enabled.
	if !k.OpenDefineAssetsChecker.IsPoolEnabled(ctx, poolId) {
		return nil, errorsmod.Wrap(types.ErrMTPDisabled, mtp.TradingAsset)
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.OpenDefineAssetsChecker.GetAmmPool(ctx, poolId, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())

	// Calculate custody amount
	custodyAmount := leveragedAmount

	switch msg.Position {
	case types.Position_LONG:
		// If collateral is not base currency, calculate the borrowing amount in base currency and check the balance
		if mtp.CollateralAsset != baseCurrency {
			custodyAmtToken := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
			borrowingAmount, err := k.OpenDefineAssetsChecker.EstimateSwapGivenOut(ctx, custodyAmtToken, baseCurrency, ammPool)
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
			custodyAmount, err = k.OpenDefineAssetsChecker.EstimateSwap(ctx, leveragedAmtTokenIn, mtp.CustodyAsset, ammPool)
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
	err = k.OpenDefineAssetsChecker.Borrow(ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, baseCurrency, isBroker)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	if err = k.OpenDefineAssetsChecker.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Take custody from the pool balance.
	if err = k.OpenDefineAssetsChecker.TakeInCustody(ctx, *mtp, &pool); err != nil {
		return nil, err
	}

	// Update the MTP health.
	lr, err := k.OpenDefineAssetsChecker.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.OpenDefineAssetsChecker.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// Update consolidated collateral amount
	k.OpenDefineAssetsChecker.CalcMTPConsolidateCollateral(ctx, mtp, baseCurrency)

	// Set stop loss price
	mtp.StopLossPrice = msg.StopLossPrice

	// Set MTP
	k.OpenDefineAssetsChecker.SetMTP(ctx, mtp)

	return mtp, nil
}
