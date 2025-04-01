package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
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

func (perpetual Perpetual) QunatityDec() math.Dec {
	return utils.IntToDec(perpetual.Quantity)
}

func (perpetual Perpetual) GetEntryValue() (math.Dec, error) {
	return perpetual.QunatityDec().Mul(perpetual.EntryPrice)
}
