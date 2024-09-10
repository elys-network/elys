package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (p PriceFeeder) GetFeederAccount() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(p.Feeder)
}
