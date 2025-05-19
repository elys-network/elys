package keeper

import (
	"errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

// GetLiquidationPrice Happens when Equity Value = Maintenance Margin Value
// Equity Value = InitialMarginValue + UPnL
// Maintenance Margin Value = MMR * Abs(Qty) * LiqPrice
// UPnL = Qty * (LiqPrice - EntryPrice)
// LiqPrice = (IMV - Qty * Entry Price) / (MMR * Abs(Qty) - Qty)
// We do not simplify this using IMR because IMR can change through gov proposal
// Also this automatically handles Liquidation price when collateral is added
func (k Keeper) GetLiquidationPrice(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket, subAccount types.SubAccount) (math.LegacyDec, error) {
	if subAccount.IsIsolated() {
		quoteDenomPrice, err := k.GetDenomPrice(ctx, market.QuoteDenom)
		if err != nil {
			return math.LegacyZeroDec(), err
		}
		initialMarginValue := perpetual.MarginAmount.ToLegacyDec().Mul(quoteDenomPrice)
		num := initialMarginValue.Sub(perpetual.Quantity.Mul(perpetual.EntryPrice))
		den := (market.MaintenanceMarginRatio.Mul(perpetual.Quantity.Abs())).Sub(perpetual.Quantity)
		// the denominator can be 0 only if MMR is 1, which would mean the position would close instantly when opened
		if den.IsZero() {
			return math.LegacyZeroDec(), errors.New("cannot calculate liquidation price, division by zero (MMR = 1)")
		}
		return num.Quo(den), nil
	} else {
		// TODO implement this for cross margin account, need research
		panic("implement me")
	}
}
