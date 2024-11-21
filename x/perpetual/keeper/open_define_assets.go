package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenDefineAssets(ctx sdk.Context, poolId uint64, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	// Determine the maximum leverage available and compute the effective leverage to be used.
	// values for leverage other than 0 or  >1 are invalidated in validate basic
	maxLeverage := k.GetMaxLeverageParam(ctx)
	proxyLeverage := sdkmath.LegacyMinDec(msg.Leverage, maxLeverage)

	// just adding collateral
	if msg.Leverage.IsZero() {
		proxyLeverage = sdkmath.LegacyOneDec()
	} else {
		// opening position, for Short we add 1 because, say atom price 5 usdc, collateral 100 usdc, leverage 5, then liabilities will be 80 atom worth 400 usdc which would be position size
		// User would be expecting position size of 100 atom / 500 usdc. So we increase the leverage from 5 to 6
		// Because of this effective leverage for short has to be reduced by 1 in query
		if msg.Position == types.Position_SHORT {
			proxyLeverage = proxyLeverage.Add(sdkmath.LegacyOneDec())
		}
		// We don't need to do this for LONG as it gives desired position
	}

	// Convert the collateral amount into a decimal format.
	collateralAmountDec := msg.Collateral.Amount.ToLegacyDec()

	// Define the assets
	var liabilitiesAsset, custodyAsset string
	switch msg.Position {
	case types.Position_LONG:
		liabilitiesAsset = baseCurrency
		custodyAsset = msg.TradingAsset
	case types.Position_SHORT:
		liabilitiesAsset = msg.TradingAsset
		custodyAsset = baseCurrency
	default:
		return nil, errorsmod.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	// Initialize a new Perpetual Trading Position (MTP).
	mtp := types.NewMTP(ctx, msg.Creator, msg.Collateral.Denom, msg.TradingAsset, liabilitiesAsset, custodyAsset, msg.Position, msg.TakeProfitPrice, poolId)

	// Call the function to process the open long logic.
	return k.ProcessOpen(ctx, mtp, proxyLeverage, collateralAmountDec, poolId, msg, baseCurrency)
}
