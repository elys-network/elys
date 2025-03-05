package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (perpetualOwner PerpetualOwner) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(perpetualOwner.Owner)
}

func (perpetual Perpetual) GetHealth(currentPrice math.LegacyDec) math.LegacyDec {
	totalValue := currentPrice.Mul(perpetual.Quantity)
	liabilities := totalValue.Sub(perpetual.Collateral.ToLegacyDec())
	h := liabilities.Quo(totalValue)
	return h
}

func (perpetual Perpetual) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(perpetual.Owner)
}

func (perpetual Perpetual) CheckEnoughMaintenence(subAccount SubAccount) bool {

	return false
}
