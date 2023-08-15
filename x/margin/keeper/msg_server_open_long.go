package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

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
