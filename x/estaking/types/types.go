package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (e ElysStaked) GetAccountAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(e.Address)
}
