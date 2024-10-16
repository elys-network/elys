package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenDefineAssets(ctx sdk.Context, poolId uint64, msg *types.MsgOpen, baseCurrency string, isBroker bool) (*types.MTP, error) {
	// Determine the maximum leverage available and compute the effective leverage to be used.
	maxLeverage := k.OpenDefineAssetsChecker.GetMaxLeverageParam(ctx)
	leverage := sdk.MinDec(msg.Leverage, maxLeverage)

	// Calculate the eta value.
	eta := leverage.Sub(sdk.OneDec())

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
	mtp := types.NewMTP(msg.Creator, msg.Collateral.Denom, msg.TradingAsset, liabilitiesAsset, custodyAsset, msg.Position, msg.TakeProfitPrice, poolId)

	// Call the function to process the open long logic.
	return k.ProcessOpen(ctx, mtp, leverage, eta, collateralAmountDec, poolId, msg, baseCurrency, isBroker)
}
