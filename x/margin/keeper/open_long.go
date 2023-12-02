package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) OpenLong(ctx sdk.Context, poolId uint64, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	// Determine the maximum leverage available and compute the effective leverage to be used.
	maxLeverage := k.OpenLongChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)

	// Calculate the eta value.
	eta := leverage.Sub(sdk.OneDec())

	// Convert the collateral amount into a decimal format.
	collateralAmountDec := sdk.NewDecFromBigInt(msg.Collateral.Amount.BigInt())

	// Define the assets
	liabilitiesAsset := ptypes.BaseCurrency
	custodyAsset := msg.TradingAsset

	// Initialize a new Margin Trading Position (MTP).
	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, liabilitiesAsset, custodyAsset, msg.Position, leverage, msg.TakeProfitPrice, poolId)

	// Call the function to process the open long logic.
	return k.ProcessOpenLong(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg, baseCurrency)
}
