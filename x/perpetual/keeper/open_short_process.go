package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) ProcessOpenShort(ctx sdk.Context, mtp *types.MTP, leverage math.LegacyDec, eta math.LegacyDec, collateralAmountDec math.LegacyDec, poolId uint64, msg *types.MsgOpen, baseCurrency string, isBroker bool) (*types.MTP, error) {
	// Fetch the pool associated with the given pool ID.
	pool, found := k.OpenShortChecker.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, mtp.TradingAsset)
	}

	// Check if the pool is enabled.
	if !k.OpenShortChecker.IsPoolEnabled(ctx, poolId) {
		return nil, errorsmod.Wrap(types.ErrMTPDisabled, mtp.TradingAsset)
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.OpenShortChecker.GetAmmPool(ctx, poolId, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := math.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())

	if mtp.CollateralAsset != baseCurrency {
		return nil, errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "collateral must be base currency")
	}

	// Define custody amount
	custodyAmount := leveragedAmount

	// check the balance
	if !types.HasSufficientPoolBalance(ammPool, mtp.CustodyAsset, custodyAmount) {
		return nil, errorsmod.Wrap(types.ErrBorrowTooHigh, custodyAmount.String())
	}

	// Borrow the asset the user wants to short.
	err = k.OpenShortChecker.Borrow(ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, baseCurrency, isBroker)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	if err = k.OpenShortChecker.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Take custody from the pool balance.
	if err = k.OpenShortChecker.TakeInCustody(ctx, *mtp, &pool); err != nil {
		return nil, err
	}

	// Update the MTP health.
	lr, err := k.OpenShortChecker.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.OpenShortChecker.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// Update consolidated collateral amount
	k.OpenShortChecker.CalcMTPConsolidateCollateral(ctx, mtp, baseCurrency)

	// Calculate consolidate liabiltiy and update consolidate leverage
	mtp.ConsolidateLeverage = types.CalcMTPConsolidateLiability(mtp)

	mtp.StopLossPrice = msg.StopLossPrice
	// Set MTP
	k.OpenShortChecker.SetMTP(ctx, mtp)

	// Return the updated Perpetual Trading Position (MTP).
	return mtp, nil
}
