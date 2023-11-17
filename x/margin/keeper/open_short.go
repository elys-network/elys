package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) OpenShort(ctx sdk.Context, poolId uint64, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	// Determine the maximum leverage available and compute the effective leverage to be used.
	maxLeverage := k.OpenShortChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)

	// Calculate the eta value.
	eta := leverage.Sub(sdk.OneDec())

	// Convert the collateral amount into a decimal format.
	collateralAmountDec := sdk.NewDecFromBigInt(msg.CollateralAmount.BigInt())

	// Initialize a new Margin Trading Position (MTP).
	mtp := types.NewMTP(msg.Creator, msg.CollateralAsset, msg.BorrowAsset, msg.Position, leverage, msg.TakeProfitPrice, poolId)

	// Call the function to process the open short logic.
	return k.ProcessOpenShort(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg, baseCurrency)
}
