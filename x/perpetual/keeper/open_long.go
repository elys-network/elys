package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenLong(ctx sdk.Context, poolId uint64, msg *types.MsgOpen, baseCurrency string, isBroker bool) (*types.MTP, error) {
	// Determine the maximum leverage available and compute the effective leverage to be used.
	maxLeverage := k.OpenLongChecker.GetMaxLeverageParam(ctx)
	leverage := sdkmath.LegacyMinDec(msg.Leverage, maxLeverage)

	// Calculate the eta value.
	eta := leverage.Sub(sdkmath.LegacyOneDec())

	// Convert the collateral amount into a decimal format.
	collateralAmountDec := sdkmath.LegacyNewDecFromBigInt(msg.Collateral.Amount.BigInt())

	// Define the assets
	liabilitiesAsset := baseCurrency
	custodyAsset := msg.TradingAsset

	// Initialize a new Perpetual Trading Position (MTP).
	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, liabilitiesAsset, custodyAsset, msg.Position, leverage, msg.TakeProfitPrice, poolId)

	// Call the function to process the open long logic.
	return k.ProcessOpenLong(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg, baseCurrency, isBroker)
}
