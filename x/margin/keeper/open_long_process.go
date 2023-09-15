package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) ProcessOpenLong(ctx sdk.Context, mtp *types.MTP, leverage sdk.Dec, eta sdk.Dec, collateralAmountDec sdk.Dec, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
	// Get token asset other than base currency
	tradingAsset := k.OpenLongChecker.GetTradingAsset(msg.CollateralAsset, msg.BorrowAsset)

	pool, found := k.OpenLongChecker.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, tradingAsset)
	}

	if !k.OpenLongChecker.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, tradingAsset)
	}

	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, poolId, tradingAsset)
	if err != nil {
		return nil, err
	}

	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())
	// If collateral is not base currency, calculate the borrowing amount in base currency and check the balance
	if msg.CollateralAsset != ptypes.BaseCurrency {
		custodyAmtToken := sdk.NewCoin(msg.CollateralAsset, leveragedAmount)
		borrowingAmount, err := k.OpenLongChecker.EstimateSwapGivenOut(ctx, custodyAmtToken, ptypes.BaseCurrency, ammPool)
		if err != nil {
			return nil, err
		}
		if !k.OpenLongChecker.HasSufficientPoolBalance(ctx, ammPool, ptypes.BaseCurrency, borrowingAmount) {
			return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
		}
	} else {
		if !k.OpenLongChecker.HasSufficientPoolBalance(ctx, ammPool, msg.CollateralAsset, leveragedAmount) {
			return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
		}
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

	// If the collateral asset is not base currency, custody amount equals to leverage amount
	if msg.CollateralAsset != ptypes.BaseCurrency {
		custodyAmount = leveragedAmount
	}

	if !k.OpenLongChecker.HasSufficientPoolBalance(ctx, ammPool, msg.BorrowAsset, custodyAmount) {
		return nil, sdkerrors.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
	}

	err = k.OpenLongChecker.Borrow(ctx, msg.CollateralAsset, msg.BorrowAsset, msg.CollateralAmount, custodyAmount, mtp, &ammPool, &pool, eta)
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

	// Update consolidated collateral amount
	k.OpenLongChecker.CalcMTPConsolidateCollateral(ctx, mtp)

	// Calculate consolidate liabiltiy
	k.OpenLongChecker.CalcMTPConsolidateLiability(ctx, mtp)

	// Set MTP
	k.OpenLongChecker.SetMTP(ctx, mtp)

	return mtp, nil
}
