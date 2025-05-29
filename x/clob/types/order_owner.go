package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (orderOwner PerpetualOrderOwner) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(orderOwner.Owner)
}
