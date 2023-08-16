package keeper

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
	paramtypes "github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.IsWhitelistingEnabled(ctx) && !k.CheckIfWhitelisted(ctx, msg.Creator) {
		return nil, sdkerrors.Wrap(types.ErrUnauthorised, "unauthorised")
	}

	if k.GetOpenMTPCount(ctx) >= (uint64)(k.GetMaxOpenPositions(ctx)) {
		return nil, sdkerrors.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}

	// If the collateral asset is not USDC
	if msg.CollateralAsset != paramtypes.USDC {
		return nil, sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "invalid collateral asset")
	}

	// Get pool Ids which can support borrowing asset
	poolIds := k.amm.GetAllPoolIdsWithDenom(ctx, msg.BorrowAsset)
	if len(poolIds) < 1 {
		return nil, sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}

	// Assume choose the first pool
	poolId := poolIds[0]

	// Get pool from pool Id
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}

	// If pool is disabled or closed in margin
	if !k.IsPoolEnabled(ctx, poolId) || k.IsPoolClosed(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, msg.BorrowAsset)
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.GetPoolOpenThreshold(ctx)) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, "pool health too low to open new positions")
	}

	var mtp *types.MTP
	var err error

	switch msg.Position {
	case types.Position_LONG:
		mtp, err = k.OpenLong(ctx, poolId, msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		sdk.NewAttribute("collateral_amount", mtp.CollateralAmount.String()),
		sdk.NewAttribute("custody_asset", mtp.CustodyAsset),
		sdk.NewAttribute("custody_amount", mtp.CustodyAmount.String()),
		sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("interest_paid_collateral", mtp.InterestPaidCollateral.String()),
		sdk.NewAttribute("interest_paid_custody", mtp.InterestPaidCustody.String()),
		sdk.NewAttribute("interest_unpaid_collateral", mtp.InterestUnpaidCollateral.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
	))

	return &types.MsgOpenResponse{}, nil
}

// Check pool balance if it is enough to borrow
func (k msgServer) CheckPoolBalance(ctx sdk.Context, ammPool ammtypes.Pool, borrowAsset string, leveragedAmount sdk.Int) bool {
	for _, asset := range ammPool.PoolAssets {
		if borrowAsset == asset.Token.Denom && asset.Token.Amount.GTE(leveragedAmount) {
			return true
		}
	}

	return false
}

// Open Long position function
func (k msgServer) OpenLong(ctx sdk.Context, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
	maxLeverage := k.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)
	eta := leverage.Sub(sdk.OneDec())

	collateralAmount := msg.CollateralAmount

	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())

	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, msg.BorrowAsset, msg.Position, leverage, poolId)

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.BorrowAsset)
	}

	if !k.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, msg.BorrowAsset)
	}

	leveragedAmountDec := collateralAmountDec.Mul(leverage)

	leveragedAmount := sdk.NewInt(leveragedAmountDec.TruncateInt().Int64())

	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.BorrowAsset)
	}

	ctx.Logger().Info(fmt.Sprintf("leveragedAmount: %s", leveragedAmount.String()))

	if !k.CheckPoolBalance(ctx, ammPool, msg.BorrowAsset, leveragedAmount) {
		return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
	}

	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, collateralAmount)
	// check if liabilities large enough for interest payments
	err := k.CheckMinLiabilities(ctx, collateralTokenAmt, eta, pool, ammPool, msg.BorrowAsset)
	if err != nil {
		return nil, err
	}

	leveragedAmtTokenIn := sdk.NewCoin(msg.CollateralAsset, leveragedAmount)
	custodyAmount, err := k.EstimateSwap(ctx, leveragedAmtTokenIn, msg.BorrowAsset, ammPool)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Info(fmt.Sprintf("custodyAmount: %s", custodyAmount.String()))

	// Custody amount check
	if !k.CheckPoolBalance(ctx, ammPool, msg.CollateralAsset, custodyAmount) {
		return nil, sdkerrors.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
	}

	err = k.Borrow(ctx, msg.CollateralAsset, collateralAmount, custodyAmount, mtp, &ammPool, &pool, eta)
	if err != nil {
		return nil, err
	}

	err = k.UpdatePoolHealth(ctx, &pool)
	if err != nil {
		return nil, err
	}

	err = k.TakeInCustody(ctx, *mtp, &pool)
	if err != nil {
		return nil, err
	}

	safetyFactor := k.GetSafetyFactor(ctx)

	lr, err := k.UpdateMTPHealth(ctx, *mtp, ammPool)

	if err != nil {
		return nil, err
	}

	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	return mtp, nil
}
