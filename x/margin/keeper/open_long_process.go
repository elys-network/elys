package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) ProcessOpenLong(ctx sdk.Context, mtp *types.MTP, leverage sdk.Dec, eta sdk.Dec, collateralAmountDec sdk.Dec, poolId uint64, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	// Fetch the pool associated with the given pool ID.
	pool, found := k.OpenLongChecker.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, mtp.TradingAsset)
	}

	// Check if the pool is enabled.
	if !k.OpenLongChecker.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, mtp.TradingAsset)
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, poolId, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())

	// If collateral is not base currency, calculate the borrowing amount in base currency and check the balance
	if mtp.CollateralAsset != baseCurrency {
		custodyAmtToken := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
		borrowingAmount, err := k.OpenLongChecker.EstimateSwapGivenOut(ctx, custodyAmtToken, baseCurrency, ammPool)
		if err != nil {
			return nil, err
		}
		if !types.HasSufficientPoolBalance(ammPool, baseCurrency, borrowingAmount) {
			return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, borrowingAmount.String())
		}
	} else {
		if !types.HasSufficientPoolBalance(ammPool, mtp.CollateralAsset, leveragedAmount) {
			return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
		}
	}

	// Check minimum liabilities.
	err = k.OpenLongChecker.CheckMinLiabilities(ctx, msg.Collateral, eta, pool, ammPool, mtp.CustodyAsset, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Calculate custody amount.
	custodyAmount := leveragedAmount
	// If position is long, calculate custody amount in custody asset
	if mtp.CollateralAsset == baseCurrency {
		leveragedAmtTokenIn := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
		custodyAmount, err = k.OpenLongChecker.EstimateSwap(ctx, leveragedAmtTokenIn, mtp.CustodyAsset, ammPool)
		if err != nil {
			return nil, err
		}
	}

	// Ensure the AMM pool has enough balance.
	if !types.HasSufficientPoolBalance(ammPool, mtp.CustodyAsset, custodyAmount) {
		return nil, sdkerrors.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
	}

	// Borrow the asset the user wants to long.
	err = k.OpenLongChecker.Borrow(ctx, msg.Collateral.Amount, custodyAmount, mtp, &ammPool, &pool, eta, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	if err = k.OpenLongChecker.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Take custody from the pool balance.
	if err = k.OpenLongChecker.TakeInCustody(ctx, *mtp, &pool); err != nil {
		return nil, err
	}

	// Update the MTP health.
	lr, err := k.OpenLongChecker.UpdateMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.OpenLongChecker.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// Update consolidated collateral amount
	k.OpenLongChecker.CalcMTPConsolidateCollateral(ctx, mtp, baseCurrency)

	// Calculate consolidate liabiltiy and update consolidate leverage
	mtp.ConsolidateLeverage = types.CalcMTPConsolidateLiability(mtp)

	// Set MTP
	k.OpenLongChecker.SetMTP(ctx, mtp)

	return mtp, nil
}
