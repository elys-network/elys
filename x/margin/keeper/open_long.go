package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) OpenLong(ctx sdk.Context, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
	// Determine the maximum leverage available and compute the effective leverage to be used.
	maxLeverage := k.OpenLongChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)

	// Calculate the eta value.
	eta := leverage.Sub(sdk.OneDec())

	// Convert the collateral amount into a decimal format.
	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())

	// Initialize a new Margin Trading Position (MTP).
	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, msg.BorrowAsset, msg.Position, leverage, poolId)

	// Call the function to process the open long logic.
	return k.ProcessOpenLong(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg)
}
