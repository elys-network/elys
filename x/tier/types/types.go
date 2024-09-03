package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (p Portfolio) GetCreatorAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(p.Creator)
}
