package types

import (
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

func NewPerpetual(id uint64, marketId uint64, owner string, qty, ep math.LegacyDec, margin math.Int, fundingRate math.LegacyDec) Perpetual {
	return Perpetual{
		Id:               id,
		MarketId:         marketId,
		Owner:            owner,
		Quantity:         qty,
		EntryPrice:       ep,
		Margin:           margin,
		EntryFundingRate: fundingRate,
	}
}
