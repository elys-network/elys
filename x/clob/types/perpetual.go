package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (perpetualOwner PerpetualOwner) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(perpetualOwner.Owner)
}

func (perpetual Perpetual) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(perpetual.Owner)
}

func (perpetual Perpetual) CheckEnoughMaintenence(subAccount SubAccount) bool {

	return false
}

func (perpetual Perpetual) IsLong() bool {
	return perpetual.Quantity.IsPositive()
}

func (perpetual Perpetual) IsShort() bool {
	return perpetual.Quantity.IsNegative()
}

func (perpetual Perpetual) IsZero() bool {
	return perpetual.Quantity.IsZero()
}

func (perpetual Perpetual) GetEntryValue() math.LegacyDec {
	return perpetual.Quantity.Abs().Mul(perpetual.EntryPrice)
}

func NewPerpetual(id uint64, marketId uint64, owner string, qty, ep math.LegacyDec, marginAmount math.Int, fundingRate math.LegacyDec) Perpetual {
	return Perpetual{
		Id:               id,
		MarketId:         marketId,
		Owner:            owner,
		Quantity:         qty,
		EntryPrice:       ep,
		MarginAmount:     marginAmount,
		EntryFundingRate: fundingRate,
	}
}

// CalculateUnrealizedPnLValue = Qty * (Current Price - EntryPrice)
func (perpetual Perpetual) CalculateUnrealizedPnLValue(markPrice math.LegacyDec) (math.LegacyDec, error) {

	if perpetual.Quantity.IsZero() || perpetual.Quantity.IsNil() {
		return math.LegacyZeroDec(), nil
	}
	if perpetual.EntryPrice.IsNil() || !perpetual.EntryPrice.IsPositive() {
		return math.LegacyDec{}, fmt.Errorf("invalid non-positive entry price: %s", perpetual.EntryPrice.String())
	}
	if markPrice.IsNil() || !markPrice.IsPositive() {
		return math.LegacyDec{}, fmt.Errorf("invalid non-positive mark price: %s", markPrice.String())
	}

	// --- PNL Calculation ---
	// Formula: UPNL = Quantity * (Mark Price - Entry Price)

	// Calculate the price difference
	priceDifference := markPrice.Sub(perpetual.EntryPrice)

	// Multiply by the quantity. The sign of the quantity automatically handles longs/shorts:
	// Long (Qty > 0): If Mark > Entry, diff is positive -> PNL is positive (profit)
	// Long (Qty > 0): If Mark < Entry, diff is negative -> PNL is negative (loss)
	// Short (Qty < 0): If Mark > Entry, diff is positive -> PNL is negative (loss)
	// Short (Qty < 0): If Mark < Entry, diff is negative -> PNL is positive (profit)
	unrealizedPNL := perpetual.Quantity.Mul(priceDifference)

	return unrealizedPNL, nil
}
