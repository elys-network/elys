package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) OpenLong(ctx sdk.Context, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
	maxLeverage := k.OpenLongChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)
	eta := leverage.Sub(sdk.OneDec())
	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())
	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, msg.BorrowAsset, msg.Position, leverage, poolId)

	// Get token asset other than USDC
	nonNativeAsset := k.GetNonNativeAsset(msg.CollateralAsset, msg.BorrowAsset)

	pool, found := k.OpenLongChecker.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, nonNativeAsset)
	}

	if !k.OpenLongChecker.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, nonNativeAsset)
	}

	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())

	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, poolId, nonNativeAsset)
	if err != nil {
		return nil, err
	}

	if !k.OpenLongChecker.HasSufficientPoolBalance(ctx, ammPool, msg.CollateralAsset, leveragedAmount) {
		return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
	}

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)

	err = k.OpenLongChecker.CheckMinLiabilities(ctx, collateralTokenAmt, eta, pool, ammPool, msg.BorrowAsset)
	if err != nil {
		return nil, err
	}

	leveragedAmtTokenIn := sdk.NewCoin(msg.CollateralAsset, leveragedAmount)
	custodyAmount, err := k.OpenLongChecker.EstimateSwap(ctx, leveragedAmtTokenIn, msg.BorrowAsset, ammPool)
	if err != nil {
		return nil, err
	}

	if !k.OpenLongChecker.HasSufficientPoolBalance(ctx, ammPool, msg.BorrowAsset, custodyAmount) {
		return nil, sdkerrors.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
	}

	err = k.OpenLongChecker.Borrow(ctx, msg.CollateralAsset, msg.CollateralAmount, custodyAmount, mtp, &ammPool, &pool, eta)
	if err != nil {
		return nil, err
	}
	if err = k.OpenLongChecker.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}
	if err = k.OpenLongChecker.TakeInCustody(ctx, *mtp, &pool); err != nil {
		return nil, err
	}

	lr, err := k.OpenLongChecker.UpdateMTPHealth(ctx, *mtp, ammPool)
	if err != nil {
		return nil, err
	}

	safetyFactor := k.OpenLongChecker.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	return mtp, nil
}
